package cyclemanager

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// func TestCycleManager(t *testing.T) {
// 	cycleInterval := 10 * time.Millisecond
// 	cycleDuration := 5 * time.Millisecond
// 	stopTimeout := 25 * time.Millisecond

// 	sleeper := sleeper{
// 		dreams: make(chan string, 1),
// 		cycleDuration: cycleDuration,
// 	}

// 	t.Run("create new", func(t *testing.T) {
// 		description := "test cycle"
// 		sleeper.sleepCycle = New(sleeper.sleep, description)

// 		assert.False(t, sleeper.sleepCycle.Running())
// 		assert.Equal(t, sleeper.sleepCycle.description, description)
// 		assert.NotNil(t, sleeper.sleepCycle.cycleFunc)
// 		assert.NotNil(t, sleeper.sleepCycle.Stop)
// 	})

// 	t.Run("start", func(t *testing.T) {
// 		sleeper.sleepCycle.Start(cycleInterval)
// 		assert.True(t, sleeper.sleepCycle.Running())
// 		assert.Equal(t, "something wonderful...", <-sleeper.dreams)
// 	})

// 	t.Run("stop", func(t *testing.T) {
// 		timeoutCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
// 		defer cancel()

// 		stopped := make(chan struct{})

// 		go func() {
// 			if sleeper.sleepCycle.TryStop(timeoutCtx) {
// 				stopped <- struct{}{}
// 			}
// 		}()

// 		select {
// 		case <-timeoutCtx.Done():
// 			t.Fatal(timeoutCtx.Err().Error(), "failed to stop sleeper")
// 		case <-stopped:
// 		}

// 		assert.False(t, sleeper.sleepCycle.Running())
// 		assert.Empty(t, <-sleeper.dreams)
// 	})
// }

// func TestCycleManager_CancelContext(t *testing.T) {
// 	cycleInterval := 10 * time.Millisecond
// 	cycleDuration := 50 * time.Millisecond
// 	stopTimeout := 25 * time.Millisecond

// 	sleeper := sleeper{
// 		dreams: make(chan string, 1),
// 		cycleDuration: cycleDuration,
// 	}

// 	t.Run("create new", func(t *testing.T) {
// 		description := "test cycle"
// 		sleeper.sleepCycle = New(sleeper.sleep, description)

// 		assert.False(t, sleeper.sleepCycle.Running())
// 		assert.Equal(t, sleeper.sleepCycle.description, description)
// 		assert.NotNil(t, sleeper.sleepCycle.cycleFunc)
// 		assert.NotNil(t, sleeper.sleepCycle.Stop)
// 	})

// 	t.Run("start", func(t *testing.T) {
// 		sleeper.sleepCycle.Start(cycleInterval)
// 		assert.True(t, sleeper.sleepCycle.Running())
// 		assert.Equal(t, "something wonderful...", <-sleeper.dreams)
// 	})

// 	t.Run("cancel early", func(t *testing.T) {
// 		ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)
// 		defer cancel()

// 		awake := make(chan struct{})

// 		go func() {
// 			if sleeper.sleepCycle.TryStop(ctx) {
// 				awake <- struct{}{}
// 			}
// 		}()

// 		select {
// 		case <-ctx.Done():
// 		case <-awake:
// 			t.Fatal("context should have been cancelled")
// 		}

// 		assert.True(t, sleeper.sleepCycle.Running())
// 		assert.Equal(t, "something wonderful...", <-sleeper.dreams)
// 		//make sure cycle was not stopped
// 		time.Sleep(2*cycleDuration)
// 		assert.True(t, sleeper.sleepCycle.Running())
// 	})
// }

// type sleeper struct {
// 	sleepCycle *CycleManager
// 	dreams     chan string
// 	cycleDuration time.Duration
// }

// func (s *sleeper) sleep(interval time.Duration) {
// 	go func() {
// 		t := time.NewTicker(interval)
// 		defer t.Stop()

// 		var ctx context.Context
// 		for {
// 			select {
// 			case ctx = <-s.sleepCycle.Stop:
// 				if (ctx.Err() == nil) {
// 					close(s.dreams)
// 					return
// 				}
// 			case <-t.C:
// 				time.Sleep(s.cycleDuration)
// 				s.dreams <- "something wonderful..."
// 			}
// 		}
// 	}()
// }
