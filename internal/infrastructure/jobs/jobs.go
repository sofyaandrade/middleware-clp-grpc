package jobs

import (
	"context"
	"sync"
)

type cancelMap struct {
	sync.Mutex
	interno map[string]context.CancelFunc
}

var (
	Jobs = newCancelMap()
)

func newCancelMap() *cancelMap {
	return &cancelMap{
		interno: make(map[string]context.CancelFunc),
	}
}

func (c *cancelMap) Get(idEquipamento string) (value context.CancelFunc, ok bool) {
	c.Lock()
	resultado, ok := c.interno[idEquipamento]
	c.Unlock()
	return resultado, ok
}

func (c *cancelMap) Set(idEquipamento string, value context.CancelFunc) {
	c.Lock()
	c.interno[idEquipamento] = value
	c.Unlock()
}

func (c *cancelMap) Delete(idEquipamento string) {
	c.Lock()
	delete(c.interno, idEquipamento)
	c.Unlock()
}
