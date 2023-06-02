package concurrency

// FanOut pattern distributes messages from single input channel to N multiply output channels.
// It's a useful pattern for parallelizing CPU and I/O utilization. Rather than coupling
// the input and computation processes in a single serial process, you might prefer to
// parallelize the workload by distributing it among some number of concurrent worker processes.
//
// FanOut is implemented as a function which accept a single source channel and integer
// representing the desired number of destination channels. FanOut creates the N destination
// channels and separate goroutines for each destination channels that compete to read the next
// value from source channel and forward to their respective destination channel.
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
