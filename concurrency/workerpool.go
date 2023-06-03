package concurrency

import (
	"sync"
)

// WorkerFunc is a function type that workers perform.
type WorkerFunc func(in any) any

// WorkerPool pattern distributes the work across multiple workers (goroutines) concurrently.
// It takes the number of workers and function of WorkerFunc type that the workers will perform.
// It also returns two channels: the first to send some value to start processing among the workers
// and the second to handle results. If the sending channel is closed, it'll close all the workers and
// WorkerPool will be finished.
func WorkerPool(n int, fn WorkerFunc) (chan<- any, <-chan any) {
	in := make(chan any)
	out := make(chan any)

	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for val := range in {
				out <- fn(val)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return in, out
}
