package clp

import (
	"context"
	"sync"
)

type Controller interface {
	Start(ctx context.Context, wg *sync.WaitGroup)
}

type Manager struct {
	controllers []Controller
}

func NewManager(controllers ...Controller) *Manager {
	return &Manager{
		controllers: controllers,
	}
}

func (m *Manager) Start(ctx context.Context, wg *sync.WaitGroup) {
	for _, controller := range m.controllers {
		wg.Add(1)
		go func(controller Controller) {
			defer wg.Done()
			controller.Start(ctx, wg)
		}(controller)
	}
}
