package cyclemanager

import (
	"context"
	"fmt"
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
	fmt.Printf("  ==> Start: started\n")
	c.Lock()
	fmt.Printf("  ==> Start: lock acquired\n")
	defer c.Unlock()

	// prevent spawning multiple cycleFunc routines
	if c.running {
		fmt.Printf("  ==> Start: already running\n")
		return
	}

	c.cycleFunc(interval)
	c.running = true
	fmt.Printf("  ==> Start: finished\n")
}

func (c *CycleManager) TryStop(ctx context.Context) (stopped bool) {
	fmt.Printf("  ==> TryStop: started\n")
	c.Lock()
	fmt.Printf("  ==> TryStop: lock acquired\n")
	defer c.Unlock()

	if !c.running {
		fmt.Printf("  ==> TryStop: already not running\n")
		return true
	}

	c.Stop <- ctx
	fmt.Printf("  ==> TryStop: stop signal read\n")
	stopped = ctx.Err() == nil
	c.running = !stopped
	fmt.Printf("  ==> TryStop: finished\n")
	return stopped
}

func (c *CycleManager) Running() bool {
	c.RLock()
	defer c.RUnlock()

	return c.running
}