package modbusmaster

import (
	"context"
	"fmt"
	"middleware/internal/domain/constants"
	"middleware/internal/domain/conversion"
	"middleware/internal/domain/interfaces"
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

const maxRegistersPerRead int = 120
const maxCoilsPerRead int = 2000
const modbusMasterTimeout = 5 * time.Second
const modbusMasterIdleTimeout = 10 * time.Second

const (
	modbusOperationCoilStatus      uint = 1
	modbusOperationInputStatus     uint = 2
	modbusOperationHoldingRegister uint = 3
	modbusOperationInputRegister   uint = 4
)

type readRegisterFunc func(address, quantity uint16) ([]byte, error)
type readBitsFunc func(address, quantity uint16) ([]byte, error)

type CommandMaster struct {
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
	clients        map[uint]modbus.Client
	handlers       map[uint]*modbus.TCPClientHandler
	channels       map[uint]chan CommandMaster
	cancelFuncs    map[uint]context.CancelFunc
	reloadRequests map[uint]struct{}
	reloadSignal   chan struct{}
}

func NewService(repository interfaces.CLPRepository) *Service {
	return &Service{
		repository:     repository,
		clients:        make(map[uint]modbus.Client),
		handlers:       make(map[uint]*modbus.TCPClientHandler),
		channels:       make(map[uint]chan CommandMaster),
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
	clps, err := s.repository.SearchClpByType(constants.MODBUS_MASTER)
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

		channel := make(chan CommandMaster, 100)
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

func (s *Service) runCLP(ctx context.Context, clp models.CLP, channel chan CommandMaster) {
	var handler *modbus.TCPClientHandler
	initializeCLPTagCache(clp)
	defer func() {
		s.unregisterCLP(clp, channel, handler)
	}()

	adress := fmt.Sprintf("%s:%s", clp.Ip, conversion.IntToString(clp.Port))
	handler = modbus.NewTCPClientHandler(adress)
	handler.SlaveId = normalizeSlaveIDMaster(clp.IdPlc)
	handler.Timeout = modbusMasterTimeout
	handler.IdleTimeout = modbusMasterIdleTimeout

RETRY:
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := ConnectWithRetry(ctx, handler)
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

		client := modbus.NewClient(handler)

		s.mu.Lock()
		s.clients[clp.ID] = client
		s.handlers[clp.ID] = handler
		s.mu.Unlock()

		ticker := time.NewTicker(1 * time.Second)
		needsRetry := false

		for {
			select {
			case cmd, ok := <-channel:
				if !ok {
					ticker.Stop()
					return
				}
				err := ExecuteCommandMaster(client, clp.ID, cmd)
				if err != nil && isConnectionError(err) {
					jobs.StatusClpRealTimeSync.Store(clp.ID, false)
					jobs.MarkCLPUnavailable(clp.ID)
					needsRetry = true
				}
			default:
			}

			if needsRetry {
				ticker.Stop()
				handler.Close()

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
				err := ExecuteCommandMaster(client, clp.ID, cmd)
				if err != nil && isConnectionError(err) {
					jobs.StatusClpRealTimeSync.Store(clp.ID, false)
					jobs.MarkCLPUnavailable(clp.ID)
					needsRetry = true
				}
			case <-ticker.C:
				select {
				case cmd, ok := <-channel:
					if !ok {
						ticker.Stop()
						return
					}
					err := ExecuteCommandMaster(client, clp.ID, cmd)
					if err != nil && isConnectionError(err) {
						jobs.StatusClpRealTimeSync.Store(clp.ID, false)
						jobs.MarkCLPUnavailable(clp.ID)
						needsRetry = true
					}
				default:
					err := updateTagsClpsMaster(client, clp.ID, clp.Tags)
					if err != nil && isConnectionError(err) {
						jobs.StatusClpRealTimeSync.Store(clp.ID, false)
						needsRetry = true
					}
				}
			}

			if !needsRetry {
				jobs.StatusClpRealTimeSync.Store(clp.ID, true)
				continue
			}

			ticker.Stop()
			handler.Close()

			select {
			case <-ctx.Done():
				return
			case <-time.After(300 * time.Millisecond):
			}

			goto RETRY
		}
	}
}

func ExecuteCommandMaster(client modbus.Client, clpId uint, cmd CommandMaster) error {
	var err error
	if cmd.Value == nil {
		readTags := cmd.Tags
		if len(readTags) == 0 && cmd.Tag != nil {
			readTags = []*models.Tag{cmd.Tag}
		}

		values, err := ReadTagsMaster(client, readTags, nil)
		if err == nil {
			cmd.Response <- values
			cmd.Erro <- nil
			return nil
		}
		time.Sleep(1 * time.Second)
	} else {
		if cmd.Tag != nil {
			err = WriteTagMaster(client, cmd.Tag, cmd.Value)
		} else {
			err = WriteTagsMaster(client, cmd.Tags, cmd.Value)
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

func updateTagsClpsMaster(client modbus.Client, clpId uint, tags []*models.Tag) error {
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

	values, err := ReadTagsMaster(client, tags, nil)
	jobs.ApplyPoll(clpId, tagIDs, values, time.Now())
	if err != nil {
		//loggar erro
		if isConnectionError(err) {
			jobs.StatusClpRealTimeSync.Store(clpId, false)
		}
	} else {
		jobs.StatusClpRealTimeSync.Store(clpId, true)
	}

	return err
}
