package utils

import (
	"context"
	"sync"
	"time"
)

// GoroutineManager is a struct that handles Goroutines with error handling and wait capabilities.
type GoroutineManager struct {
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error
	timeout time.Duration
	mu      sync.Mutex
}

// NewGoroutineManager initializes a new GoroutineManager with a timeout for each Goroutine.
func NewGoroutineManager(timeout time.Duration) *GoroutineManager {
	return &GoroutineManager{
		timeout: timeout,
	}
}

// Go starts a new Goroutine and tracks any error returned by the function.
func (gm *GoroutineManager) Go(f func(ctx context.Context) error) {
	gm.wg.Add(1)
	go func() {
		defer gm.wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), gm.timeout)
		defer cancel()

		if err := f(ctx); err != nil {
			gm.mu.Lock()
			gm.errOnce.Do(func() {
				gm.err = err
			})
			gm.mu.Unlock()
		}
	}()
}

// Wait waits for all Goroutines to complete. Returns the first error encountered, if any.
func (gm *GoroutineManager) Wait() error {
	gm.wg.Wait()
	return gm.err
}
