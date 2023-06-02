package concurrency

import "context"

// Generator is used to generate the sequence of values. It simply returns a channel
// from which we can read the values. This is a similar behavior as yield in JavaScript and Python.
func Generator(ctx context.Context) <-chan int {
	ch := make(chan int)

	go func() {
		for i := 0; ; i++ {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case ch <- i:
			}
		}
	}()

	return ch
}
