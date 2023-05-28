package concurrency

import "sync"

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
