package concurrency

import "testing"

func TestFanIn(t *testing.T) {
	sources := make([]<-chan int, 3)
	for i := range sources {
		src := make(chan int)
		sources[i] = src

		go func() {
			for j := 1; j <= 5; j++ {
				src <- j
			}

			close(src)
		}()
	}

	var sum int
	dest := FanIn(sources...)

	for val := range dest {
		sum += val
	}

	if sum != 45 {
		t.Errorf("wrong FanIn result: got %d, expected 45", sum)
	}
}
