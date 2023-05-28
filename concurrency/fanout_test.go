package concurrency

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestFanOut(t *testing.T) {
	const n = 5
	var expectedTotalSum int64

	src := make(chan int)

	go func() {
		var sum int

		for i := 1; i <= 10; i++ {
			src <- i
			sum += i
		}

		atomic.AddInt64(&expectedTotalSum, int64(sum))
		close(src)
	}()

	var wg sync.WaitGroup
	var totalSum int64

	dests := FanOut(src, n)
	wg.Add(n)

	for _, dest := range dests {
		go func(dest <-chan int) {
			defer wg.Done()

			var sum int
			for val := range dest {
				sum += val
			}

			atomic.AddInt64(&totalSum, int64(sum))
		}(dest)
	}

	wg.Wait()

	if totalSum != expectedTotalSum {
		t.Errorf("wrong fan out result: got %d, want %d", totalSum, expectedTotalSum)
	}
}
