package cyclemanager

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type CycleFunc func()

type CycleManager2 struct {
	sync.RWMutex

	cycleFunc     CycleFunc
	cycleInterval time.Duration
	running       bool
	stop          chan *stopSignal
}

type stopSignal struct {
	contexts   []context.Context
	stopResult chan bool
}

func New2(cycleInterval time.Duration, cycleFunc CycleFunc) *CycleManager2 {
	return &CycleManager2{
		cycleFunc:     cycleFunc,
		cycleInterval: cycleInterval,
		running:       false,
		stop:          make(chan *stopSignal, 1),
	}
}

func (c *CycleManager2) Start() {
	fmt.Printf("   ==> Start: beginning\n")
	c.Lock()
	fmt.Printf("   ==> Start: lock acquired\n")
	defer c.Unlock()

	if c.running {
		fmt.Printf("   ==> Start: already running\n")
		return
	}

	handleStopSignal := func(sig *stopSignal) bool {
		defer close(sig.stopResult)

		fmt.Printf("   ==> loop: signal start\n")
		proceed := false
		for _, ctx := range sig.contexts {
			if ctx.Err() == nil {
				proceed = true
				break
			}
		}
		if proceed {
			c.Lock()
			c.running = false
			c.Unlock()

			sig.stopResult <- true
			fmt.Printf("   ==> loop: signal end (stopped)\n")
			return false
		}
		sig.stopResult <- false
		fmt.Printf("   ==> loop: signal end (cancelled)\n")

		return true
	}

	go func() {
		fmt.Printf("   ==> loop: created\n")

		ticker := time.NewTicker(c.cycleInterval)
		defer ticker.Stop()

		var sig *stopSignal
		for {
			fmt.Printf("   ==> loop: beginning, chan len [%+v]\n", len(c.stop))
			select {
			case sig = <-c.stop:
				fmt.Printf("   ==> loop: in stop case\n")
				if ok := handleStopSignal(sig); !ok {
					return
				}
			case <-ticker.C:
				select {
				case sig = <-c.stop:
					fmt.Printf("   ==> loop: in tick case, but stop exists\n")
					if ok := handleStopSignal(sig); !ok {
						return
					}
				default:
				}

				fmt.Printf("   ==> loop: tick start\n")
				c.cycleFunc()
				fmt.Printf("   ==> loop: tick end\n")
			}
			fmt.Printf("   ==> loop: finish\n")
		}
	}()

	c.running = true
	fmt.Printf("   ==> Start: finish\n")
}

func (c *CycleManager2) TryStop(ctx context.Context) (stopResult chan bool) {
	fmt.Printf("   ==> TryStop: beginning\n")
	c.Lock()
	fmt.Printf("   ==> TryStop: lock acquired\n")
	defer c.Unlock()

	stopResult = make(chan bool, 1)
	if !c.running {
		stopResult <- true
		close(stopResult)
		return stopResult
	}

	select {
	// there is already pending stop, add new context to previous ones
	case prevSignal := <-c.stop:
		fmt.Printf("   ==> TryStop: prev signal\n")
		commonStopped := make(chan bool, 1)
		c.stop <- &stopSignal{append(prevSignal.contexts, ctx), commonStopped}
		go func() {
			cs := <-commonStopped
			prevSignal.stopResult <- cs
			close(prevSignal.stopResult)
			stopResult <- cs
			close(stopResult)
		}()
	default:
		fmt.Printf("   ==> TryStop: default\n")
		c.stop <- &stopSignal{[]context.Context{ctx}, stopResult}
		fmt.Printf("   ==> TryStop: default2\n")
	}

	fmt.Printf("   ==> TryStop: finish\n")
	return stopResult
}

func (c *CycleManager2) Running() bool {
	c.RLock()
	defer c.RUnlock()

	return c.running
}
