package modbusmaster

import (
	"middleware/internal/domain/models"
	"middleware/internal/infrastructure/jobs"

	"github.com/goburrow/modbus"
)

func (s *Service) unregisterCLP(clp models.CLP, channel chan CommandMaster, handler *modbus.TCPClientHandler) {
	if handler != nil {
		handler.Close()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	currentChannel, ok := s.channels[clp.ID]
	if ok && currentChannel != channel {
		return
	}

	delete(s.clients, clp.ID)
	delete(s.handlers, clp.ID)
	delete(s.channels, clp.ID)
	delete(s.quantityTags, clp.ID)
	delete(s.signaturesTags, clp.ID)
	delete(s.cancelFuncs, clp.ID)
	delete(s.settings, clp.ID)

	jobs.MutexMaster.Lock()
	delete(jobs.TagsByClpMaster, clp.ID)
	jobs.MutexMaster.Unlock()

	for _, tag := range clp.Tags {
		if tag != nil {
			jobs.TagsByIDMaster.Delete(tag.ID)
		}
	}
	jobs.StatusClpRealTimeSync.Delete(clp.ID)
}

func (s *Service) stopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for clpID, cancel := range s.cancelFuncs {
		cancel()
		delete(s.cancelFuncs, clpID)
	}

	for clpID, channel := range s.channels {
		close(channel)
		delete(s.channels, clpID)
	}
}

func (s *Service) stopMissingCLPs(activeCLPs map[uint]struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for clpID, cancel := range s.cancelFuncs {
		if _, ok := activeCLPs[clpID]; ok {
			continue
		}

		cancel()
		delete(s.cancelFuncs, clpID)

		if channel, ok := s.channels[clpID]; ok {
			close(channel)
			delete(s.channels, clpID)
		}
	}
}
