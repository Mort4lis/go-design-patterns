package concurrency

import (
	"runtime"
	"testing"
	"time"
)

func incrJob(in <-chan any, out chan<- any) {
	for val := range in {
		v := val.(int)
		v++

		out <- v
	}
}

func TestPipeline(t *testing.T) {
	var want int
	jobs := make([]Job, 5)

	for i := range jobs {
		want += 1
		jobs[i] = incrJob
	}

	beforeNumGs := runtime.NumGoroutine()

	in, out := Pipeline(jobs...)
	in <- 0
	val := <-out

	if val != want {
		t.Errorf("wrong pipeline result: got %d, want %d", val, want)
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

	afterNumGs := runtime.NumGoroutine()
	if beforeNumGs < afterNumGs {
		buf := make([]byte, 4096)
		runtime.Stack(buf, true)

		t.Log(string(buf))
		t.Errorf(
			"gorotines leek is detected: before running pipeline was %d, now %d",
			beforeNumGs, afterNumGs,
		)
	}
}
