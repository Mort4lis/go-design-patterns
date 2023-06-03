package concurrency

import (
	"testing"
	"time"
)

func longSquare(in any) any {
	val := in.(int)
	time.Sleep(2 * time.Second)

	return val * val
}

func TestWorkerPool(t *testing.T) {
	const numWorkers = 5
	var (
		want int
		got  int
	)

	in, out := WorkerPool(numWorkers, longSquare)
	time.Sleep(50 * time.Millisecond)

	for i := 1; i <= numWorkers; i++ {
		select {
		case in <- i:
		default:
			t.Fatal("no one reads an input channel")
		}

		want += i * i
	}

	timer := time.NewTimer(2200 * time.Millisecond)

	for i := 1; i <= numWorkers; i++ {
		var val any

		select {
		case val = <-out:
		case <-timer.C:
			t.Fatal("timeout exceeded: no one writes an output channel")
		}

		got += val.(int)
	}

	close(in)
	time.Sleep(50 * time.Millisecond)

	isOutClosed := false

	select {
	case _, ok := <-out:
		isOutClosed = !ok
	default:
	}

	if !isOutClosed {
		t.Errorf("out channel isn't closed, but expected")
	}

	if got != want {
		t.Errorf("wrong worker pool result: got %d, want %d", got, want)
	}
}
