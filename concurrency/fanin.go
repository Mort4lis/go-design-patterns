package concurrency

import "sync"

// FanIn pattern combines multiply inputs into one single output channel.
// Services that have some number of workers that all generate output may find it useful
// to combine all the workers’ outputs to be processed as a single unified stream.
//
// FanIn is implemented as a function that receives N source channels. For each
// input channel FanIn starts a separate goroutine to read values from its assigned channel and
// forward all the values to a single destination channel shared by all the goroutines.
func FanIn(sources ...<-chan int) <-chan int {
	dest := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(sources))

	for _, src := range sources {
		go func(ch <-chan int) {
			defer wg.Done()

			for val := range ch {
				dest <- val
			}
		}(src)
	}

	go func() {
		wg.Wait()
		close(dest)
	}()

	return dest
}
