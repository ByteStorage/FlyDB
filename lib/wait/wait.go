package wait

import (
	"sync"
	"time"
)

// Wait is a wrapper of sync.WaitGroup which can wait with timeout.
type Wait struct {
	wg sync.WaitGroup
}

// NewWait returns a new Wait.
func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

// Done decrements the WaitGroup counter by one.
func (w *Wait) Done() {
	w.wg.Done()
}

// Wait waits until the WaitGroup counter is zero.
func (w *Wait) Wait() {
	w.wg.Wait()
}

// WaitTimeout waits until the WaitGroup counter is zero or timeout.
// It returns true if timed out.
func (w *Wait) WaitTimeout(timeout time.Duration) bool {
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- true
	}()
	select {
	case <-c:
		return false
	case <-time.After(timeout): // complete normally
		return true // timed out
	}
}
