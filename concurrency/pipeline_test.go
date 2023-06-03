package concurrency

import (
	"testing"
	"time"

	"github.com/Mort4lis/go-design-patterns/testutils"
)

func incrJob(in <-chan any, out chan<- any) {
	for val := range in {
		v := val.(int)
		v++

		out <- v
	}
}

func TestPipeline(t *testing.T) {
	testutils.DetectGoroutineLeeks(t, func() {
		var want int
		jobs := make([]Job, 5)

		for i := range jobs {
			want += 1
			jobs[i] = incrJob
		}

		in, out := Pipeline(jobs...)

		testutils.WriteChan(t, in, 0, testutils.WithDuration(50*time.Millisecond))
		val := testutils.ReadChan(t, out, testutils.WithDuration(50*time.Millisecond))

		if val != want {
			t.Errorf("wrong pipeline result: got %d, want %d", val, want)
		}

		close(in)
		testutils.CheckClosedChan(t, out, testutils.WithDuration(50*time.Millisecond))
	})
}
