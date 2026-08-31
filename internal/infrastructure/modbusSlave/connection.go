package modbusSlave

import (
	"context"
	"middleware/internal/domain/conversion"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
	"net"
	"sync"
	"time"
)

const maxRegistersPerRead int = 120
const maxCoilsPerRead int = 2000
const modbusSlaveTimeout = 5 * time.Second
const modbusSlaveIdleTimeout = 10 * time.Second

const (
	modbusOperationCoilStatus      uint = 1
	modbusOperationInputStatus     uint = 2
	modbusOperationHoldingRegister uint = 3
	modbusOperationInputRegister   uint = 4
)

type readRegisterFunc func(address, quantity uint16) ([]byte, error)
type readBitsFunc func(address, quantity uint16) ([]byte, error)

type CommandSlave struct {
	Tag      *models.Tag
	Tags     []*models.Tag
	Value    interface{}
	Response chan interface{}
	Erro     chan error
}

type readBlock struct {
	start        int
	endExclusive int
	tags         []*models.Tag
}

type Service struct {
	repository interfaces.CLPRepository

	mu             sync.Mutex
	clients        map[uint]modbusSlaveClient
	handlers       map[uint]*slaveTCPServer
	channels       map[uint]chan CommandSlave
	cancelFuncs    map[uint]context.CancelFunc
	reloadRequests map[uint]struct{}
	reloadSignal   chan struct{}
}

func NewService(repository interfaces.CLPRepository) *Service {
	return &Service{
		repository:     repository,
		clients:        make(map[uint]modbusSlaveClient),
		handlers:       make(map[uint]*slaveTCPServer),
		channels:       make(map[uint]chan CommandSlave),
		cancelFuncs:    make(map[uint]context.CancelFunc),
		reloadRequests: make(map[uint]struct{}),
		reloadSignal:   make(chan struct{}, 1),
	}
}

func (s *Service) RequestCLPReload(clpID uint) {
	if clpID == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.reloadRequests[clpID] = struct{}{}

	select {
	case s.reloadSignal <- struct{}{}:
	default:
	}
}

func (s *Service) Start(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	s.syncCLPs(ctx, wg)

	for {
		select {
		case <-ctx.Done():
			s.stopAll()
			return
		case <-s.reloadSignal:
			s.syncCLPs(ctx, wg)
		case <-ticker.C:
			s.syncCLPs(ctx, wg)
		}
	}
}

func (s *Service) syncCLPs(ctx context.Context, wg *sync.WaitGroup) {
	clps, err := s.repository.SearchClpByType(2)
	if err != nil {
		return
	}

	activeCLPs := make(map[uint]struct{}, len(*clps))
	for _, clp := range *clps {
		activeCLPs[clp.ID] = struct{}{}
		s.reloadCLPIfNeeded(ctx, wg, clp)
	}
	s.stopMissingCLPs(activeCLPs)
}

func (s *Service) reloadCLPIfNeeded(ctx context.Context, wg *sync.WaitGroup, clp models.CLP) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, running := s.channels[clp.ID]
	_, requestedReload := s.reloadRequests[clp.ID]
	needsReload := !running || requestedReload

	if needsReload {
		delete(s.reloadRequests, clp.ID)

		if cancel, ok := s.cancelFuncs[clp.ID]; ok {
			cancel()
			delete(s.cancelFuncs, clp.ID)
		}

		if channel, ok := s.channels[clp.ID]; ok {
			close(channel)
			delete(s.channels, clp.ID)
		}

		channel := make(chan CommandSlave, 100)
		s.channels[clp.ID] = channel

		clpCtx, cancel := context.WithCancel(ctx)
		s.cancelFuncs[clp.ID] = cancel

		wg.Add(1)
		go func() {
			defer wg.Done()
			s.runCLP(clpCtx, clp, channel)
		}()
	}
}

func (s *Service) runCLP(ctx context.Context, clp models.CLP, channel chan CommandSlave) {
	var server *slaveTCPServer
	initializeCLPTagCache(clp)
	defer func() {
		s.unregisterCLP(clp, channel, server)
	}()

	address := net.JoinHostPort(clp.Ip, conversion.IntToString(clp.Port))
	server = newSlaveTCPServer(address, normalizeSlaveIDSlave(clp.IdPlc), clp.ID, clp.Tags, modbusSlaveTimeout)

RETRY:
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := ConnectWithRetry(ctx, server)
		if err != nil {
			jobs.StatusClpRealTimeSync.Store(clp.ID, false)
			jobs.MarkCLPUnavailable(clp.ID)

			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				continue RETRY
			}
		}

		s.mu.Lock()
		s.clients[clp.ID] = server
		s.handlers[clp.ID] = server
		s.mu.Unlock()
		server.PublishTags()
		jobs.StatusClpRealTimeSync.Store(clp.ID, true)

		ticker := time.NewTicker(1 * time.Second)
		needsRetry := false

		for {
			select {
			case cmd, ok := <-channel:
				if !ok {
					ticker.Stop()
					return
				}
				err := ExecuteCommandSlave(server, clp.ID, cmd)
				if err != nil && isConnectionError(err) {
					jobs.StatusClpRealTimeSync.Store(clp.ID, false)
					jobs.MarkCLPUnavailable(clp.ID)
					needsRetry = true
				}
			default:
			}

			if needsRetry {
				ticker.Stop()
				server.Close()

				select {
				case <-ctx.Done():
					return
				case <-time.After(300 * time.Millisecond):
				}

				goto RETRY
			}

			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case cmd, ok := <-channel:
				if !ok {
					ticker.Stop()
					return
				}
				err := ExecuteCommandSlave(server, clp.ID, cmd)
				if err != nil && isConnectionError(err) {
					jobs.StatusClpRealTimeSync.Store(clp.ID, false)
					jobs.MarkCLPUnavailable(clp.ID)
					needsRetry = true
				}
			case <-server.Errors():
				jobs.StatusClpRealTimeSync.Store(clp.ID, false)
				jobs.MarkCLPUnavailable(clp.ID)
				needsRetry = true
			case <-ticker.C:
				select {
				case cmd, ok := <-channel:
					if !ok {
						ticker.Stop()
						return
					}
					err := ExecuteCommandSlave(server, clp.ID, cmd)
					if err != nil && isConnectionError(err) {
						jobs.StatusClpRealTimeSync.Store(clp.ID, false)
						jobs.MarkCLPUnavailable(clp.ID)
						needsRetry = true
					}
				default:
					server.PublishTags()
				}
			}

			if !needsRetry {
				jobs.StatusClpRealTimeSync.Store(clp.ID, true)
				continue
			}

			ticker.Stop()
			server.Close()

			select {
			case <-ctx.Done():
				return
			case <-time.After(300 * time.Millisecond):
			}

			goto RETRY
		}
	}
}

func ExecuteCommandSlave(client modbusSlaveClient, clpId uint, cmd CommandSlave) error {
	var err error
	if cmd.Value == nil {
		readTags := cmd.Tags
		if len(readTags) == 0 && cmd.Tag != nil {
			readTags = []*models.Tag{cmd.Tag}
		}

		values, err := ReadTagsSlave(client, readTags, nil)
		if err == nil {
			cmd.Response <- values
			cmd.Erro <- nil
			return nil
		}
		time.Sleep(1 * time.Second)
	} else {
		if cmd.Tag != nil {
			err = WriteTagSlave(client, cmd.Tag, cmd.Value)
		} else {
			err = WriteTagsSlave(client, cmd.Tags, cmd.Value)
		}
		if err == nil {
			cmd.Response <- nil
			cmd.Erro <- nil
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	if isConnectionError(err) {
		return err
	}

	_ = clpId
	cmd.Response <- nil
	cmd.Erro <- err
	return err
}

func updateTagsClpsSlave(client modbusSlaveClient, clpId uint, tags []*models.Tag) error {
	tagIDs := make([]uint, 0, len(tags))
	for _, tag := range tags {
		if tag != nil {
			tagIDs = append(tagIDs, tag.ID)
		}
	}

	if len(tags) == 0 {
		jobs.ApplyPoll(clpId, tagIDs, nil, time.Now())
		return nil
	}

	values, err := ReadTagsSlave(client, tags, nil)
	jobs.ApplyPoll(clpId, tagIDs, values, time.Now())
	if err != nil {
		if isConnectionError(err) {
			jobs.StatusClpRealTimeSync.Store(clpId, false)
		}
	} else {
		jobs.StatusClpRealTimeSync.Store(clpId, true)
	}

	return err
}
