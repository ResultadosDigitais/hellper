package concurrence

import "sync"

func WithWaitGroup(wg *sync.WaitGroup, fn func()) {
	wg.Add(1)
	go func() {
		fn()
		wg.Done()
	}()
}
