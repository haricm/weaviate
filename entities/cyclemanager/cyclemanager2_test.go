package cyclemanager

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCycleManager2(t *testing.T) {
	cycleInterval := 5 * time.Millisecond
	cycleDuration := 2 * time.Millisecond
	stopTimeout := 12 * time.Millisecond

	results := make(chan string, 1)
	var cm *CycleManager2

	t.Run("create new", func(t *testing.T) {
		cm = New2(cycleInterval, func() {
			time.Sleep(cycleDuration)
			results <- "something wonderful..."
		})

		assert.False(t, cm.Running())
		assert.Equal(t, cycleInterval, cm.cycleInterval)
		assert.NotNil(t, cm.cycleFunc)
		assert.NotNil(t, cm.stop)
	})

	t.Run("start", func(t *testing.T) {
		cm.Start()

		assert.True(t, cm.Running())
		assert.Equal(t, "something wonderful...", <-results)
	})

	t.Run("stop", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel()

		stopResult := cm.TryStop(timeoutCtx)
		fmt.Printf("   ==> test/after trystop [%+v]\n", stopResult)

		select {
		case <-timeoutCtx.Done():
			fmt.Printf("   ==> timeout\n")
			t.Fatal(timeoutCtx.Err().Error(), "failed to stop")
		case stopped := <-stopResult:
			fmt.Printf("   ==> stopResult\n")
			assert.True(t, stopped)
			assert.False(t, cm.Running())
			assert.Empty(t, results)
		}
		fmt.Printf("   ==> after select\n")
	})

	fmt.Printf("   ==> all end\n")
}

func TestCycleManager2_timeout(t *testing.T) {
	cycleInterval := 5 * time.Millisecond
	cycleDuration := 20 * time.Millisecond
	stopTimeout := 12 * time.Millisecond

	results := make(chan string, 1)
	cm := New2(cycleInterval, func() {
		time.Sleep(cycleDuration)
		fmt.Printf("   ==> slept\n")
		results <- "something wonderful..."
	})
	
	t.Run("timeout is reached", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel()
	
		cm.Start()
	
		// wait for 1st cycle to start
		time.Sleep(cycleInterval + 1 * time.Millisecond)
		stopResult := cm.TryStop(timeoutCtx)
		fmt.Printf("   ==> test/after trystop [%+v]\n", stopResult)

		select {
		case <-timeoutCtx.Done():
			fmt.Printf("   ==> timeout\n")
			assert.True(t, cm.Running())
		case <-stopResult:
			t.Fatal("stopped before timeout")
		}
		fmt.Printf("   ==> after select\n")

		// make sure it is still running
		assert.False(t, <-stopResult)
		assert.True(t, cm.Running())
		assert.Equal(t, "something wonderful...", <-results)
	})

	t.Run("stop", func(t *testing.T) {
		stopResult := cm.TryStop(context.Background())
		assert.True(t, <-stopResult)
		assert.False(t, cm.Running())
	})

	fmt.Printf("   ==> all end\n")
}

func TestCycleManager2_doesNotStartMultipleTimes(t *testing.T) {
	cycleInterval := 5 * time.Millisecond
	cycleDuration := 2 * time.Millisecond
	
	startCount := 5

	results := make(chan string, startCount)
	cm := New2(cycleInterval, func() {
		time.Sleep(cycleDuration)
		results <- "something wonderful..."
	})

	t.Run("multiple starts", func(t *testing.T) {
		for i := 0; i < startCount; i++ {
			cm.Start()
		}

		// wait for 1st cycle to start
		time.Sleep(cycleInterval + 1 * time.Millisecond)
		stopResult := cm.TryStop(context.Background())

		assert.True(t, <-stopResult)
		assert.False(t, cm.Running())
		// just one result produced
		assert.Equal(t, 1, len(results))
	})
}

func TestCycleManager2_handlesMultipleStops(t *testing.T) {
	cycleInterval := 5 * time.Millisecond
	cycleDuration := 2 * time.Millisecond
	
	stopCount := 5

	results := make(chan string, 1)
	cm := New2(cycleInterval, func() {
		time.Sleep(cycleDuration)
		results <- "something wonderful..."
	})

	t.Run("multiple stops", func(t *testing.T) {
		cm.Start()

		// wait for 1st cycle to start
		time.Sleep(cycleInterval + 1 * time.Millisecond)
		stopResult := make([]chan bool, stopCount)
		for i := 0; i < stopCount; i++ {
			stopResult[i] = cm.TryStop(context.Background())
		}
		
		for i := 0; i < stopCount; i++ {
			assert.True(t, <-stopResult[i])
		}

		assert.False(t, cm.Running())
		assert.Equal(t, "something wonderful...", <-results)
	})
}

func TestCycleManager2_stopsIfNotAllContextsAreCancelled(t *testing.T) {
	cycleInterval := 5 * time.Millisecond
	cycleDuration := 2 * time.Millisecond
	stopTimeout := 5 * time.Millisecond
			
	results := make(chan string, 1)
	cm := New2(cycleInterval, func() {
		time.Sleep(cycleDuration)
		results <- "something wonderful..."
	})

	t.Run("multiple stops, few cancelled", func(t *testing.T) {
		timeout1Ctx, cancel1 := context.WithTimeout(context.Background(), stopTimeout)
		timeout2Ctx, cancel2 := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel1()
		defer cancel2()
	
		cm.Start()

		// wait for 1st cycle to start
		time.Sleep(cycleInterval + 1 * time.Millisecond)

		stopResult1 := cm.TryStop(timeout1Ctx)
		stopResult2 := cm.TryStop(timeout2Ctx)
		stopResult3 := cm.TryStop(context.Background())

		// all produce the same result: cycle was stopped
		assert.True(t, <-stopResult1)
		assert.True(t, <-stopResult2)
		assert.True(t, <-stopResult3)

		assert.False(t, cm.Running())
		assert.Equal(t, "something wonderful...", <-results)
	})
}

func TestCycleManager2_doesNotStopIfAllContextsAreCancelled(t *testing.T) {
	cycleInterval := 5 * time.Millisecond
	cycleDuration := 2 * time.Millisecond
	stopTimeout := 5 * time.Millisecond
			
	results := make(chan string, 1)
	cm := New2(cycleInterval, func() {
		time.Sleep(cycleDuration)
		results <- "something wonderful..."
	})

	t.Run("multiple stops, few cancelled", func(t *testing.T) {
		timeout1Ctx, cancel1 := context.WithTimeout(context.Background(), stopTimeout)
		timeout2Ctx, cancel2 := context.WithTimeout(context.Background(), stopTimeout)
		timeout3Ctx, cancel3 := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel1()
		defer cancel2()
		defer cancel3()
	
		cm.Start()

		// wait for 1st cycle to start
		time.Sleep(cycleInterval + 1 * time.Millisecond)

		stopResult1 := cm.TryStop(timeout1Ctx)
		stopResult2 := cm.TryStop(timeout2Ctx)
		stopResult3 := cm.TryStop(timeout3Ctx)

		// all produce the same result: cycle was stopped
		assert.False(t, <-stopResult1)
		assert.False(t, <-stopResult2)
		assert.False(t, <-stopResult3)

		assert.True(t, cm.Running())
		assert.Equal(t, "something wonderful...", <-results)
	})

	t.Run("stop", func(t *testing.T) {
		stopResult := cm.TryStop(context.Background())
		assert.True(t, <-stopResult)
		assert.False(t, cm.Running())
	})
}
