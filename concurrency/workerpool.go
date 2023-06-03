package concurrency

import (
	"sync"
)

type WorkerFunc func(in any) any

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
