package concurrency

import (
	"sync/atomic"
	"testing"
)

func TestFanIn(t *testing.T) {
	var expectedTotalSum int64
	sources := make([]<-chan int, 3)

	for i := range sources {
		src := make(chan int)
		sources[i] = src

		go func() {
			var sum int

			for j := 1; j <= 5; j++ {
				src <- j
				sum += j
			}

			atomic.AddInt64(&expectedTotalSum, int64(sum))
			close(src)
		}()
	}

	var totalSum int
	dest := FanIn(sources...)

	for val := range dest {
		totalSum += val
	}

	if totalSum != int(expectedTotalSum) {
		t.Errorf("wrong FanIn result: got %d, want %d", totalSum, expectedTotalSum)
	}
}
