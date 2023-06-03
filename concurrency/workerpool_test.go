package concurrency

import (
	"testing"
	"time"

	"github.com/Mort4lis/go-design-patterns/testutils"
)

func longSquare(in any) any {
	val := in.(int)
	time.Sleep(2 * time.Second)

	return val * val
}

func TestWorkerPool(t *testing.T) {
	testutils.DetectGoroutineLeeks(t, func() {
		const numWorkers = 5
		var (
			want int
			got  int
		)

		in, out := WorkerPool(numWorkers, longSquare)
		writeTimer := time.NewTimer(100 * time.Millisecond)

		for i := 1; i <= numWorkers; i++ {
			testutils.WriteChan(t, in, i, testutils.WithTimer(writeTimer))
			want += i * i
		}

		readTimer := time.NewTimer(2100 * time.Millisecond)

		for i := 1; i <= numWorkers; i++ {
			val := testutils.ReadChan(t, out, testutils.WithTimer(readTimer))
			got += val.(int)
		}

		close(in)
		testutils.CheckClosedChan(t, out, testutils.WithDuration(50*time.Millisecond))

		if got != want {
			t.Errorf("wrong worker pool result: got %d, want %d", got, want)
		}
	})
}
