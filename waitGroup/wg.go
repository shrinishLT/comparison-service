package waitGroup

import (
	"sync"
	"time"
)

var wg sync.WaitGroup

// GetGlobalWaitGroup returns global wait group
func GetGlobalWaitGroup() *sync.WaitGroup {
	return &wg
}

// WaitWithTimeout wait blocks until the WaitGroup counter is zero but unblocks if the timeout is reached. Returns true
// if the wait group completes before timeout.
func WaitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}
