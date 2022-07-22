package cyclemanager

import (
	"context"
	"sync"
	"time"
)

type CycleManager struct {
	sync.RWMutex

	description string
	running     bool
	cycleFunc   func(duration time.Duration)

	Stop chan context.Context
}

func New(cycleFunc func(duration time.Duration), description string) *CycleManager {
	return &CycleManager{
		description: description,
		cycleFunc:   cycleFunc,
		Stop:        make(chan context.Context),
	}
}

func (c *CycleManager) Start(interval time.Duration) {
	c.Lock()
	defer c.Unlock()

	// prevent spawning multiple cycleFunc routines
	if c.running {
		return
	}

	c.cycleFunc(interval)
	c.running = true
}

func (c *CycleManager) TryStop(ctx context.Context) (stopped bool) {
	c.Lock()
	defer c.Unlock()

	if !c.running {
		return true
	}

	c.Stop <- ctx
	stopped = ctx.Err() == nil
	c.running = !stopped
	return stopped
}

func (c *CycleManager) Running() bool {
	c.RLock()
	defer c.RUnlock()

	return c.running
}