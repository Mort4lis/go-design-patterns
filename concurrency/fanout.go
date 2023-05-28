package concurrency

func FanOut(src <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, n)

	for i := range dests {
		dest := make(chan int)
		dests[i] = dest

		go func() {
			for val := range src {
				dest <- val
			}

			close(dest)
		}()
	}

	return dests
}
