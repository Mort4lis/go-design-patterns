package concurrency

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

func slowFunction() (string, error) {
	time.Sleep(2 * time.Second)

	return "I slept for 2 seconds", nil
}

// TestFuture just runs slow functions, and makes sure that it returns the
// expected result after the expected amount of time.
func TestFuture(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := RunAsync(ctx, slowFunction)

	res, err := future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 2)
}

// TestFutureGetTwice tests that subsequent calls to future.Result()
// immediately return the initial return values.
func TestFutureGetTwice(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := RunAsync(ctx, slowFunction)

	res, err := future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 2)

	// Get result again. Should happen straightaway.

	start = time.Now()

	res, err = future.Result()
	if err != nil {
		t.Error(err)
		return
	}

	if !strings.HasPrefix(res, "I slept for") {
		t.Error("unexpected output:", res)
	}

	elapsedCheck(t, start, 0)
}

// TestFutureConcurrent tests that the Future is thread-safe.
func TestFutureConcurrent(t *testing.T) {
	start := time.Now()

	ctx := context.Background()
	future := RunAsync(ctx, slowFunction)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			res, err := future.Result()
			if err != nil {
				t.Error(err)
				return
			}

			if !strings.HasPrefix(res, "I slept for") {
				t.Error("unexpected output:", res)
			}

			elapsedCheck(t, start, 2)
		}()
	}

	wg.Wait()
}

// TestFutureTimeout makes sure that the future will time out with an error
// if its context is canceled with a timeout
func TestFutureTimeout(t *testing.T) {
	start := time.Now()

	// Get a context decorated with a 1-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	// The cancel function returned by context.WithTimeout should be called,
	// not discarded, to avoid a context leak
	defer cancel()

	future := RunAsync(ctx, slowFunction)

	// We should time out with a "context deadline exceeded" error
	res, err := future.Result()
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Error("received unexpected error: ", err)
	}

	// Result should be empty
	if res != "" {
		t.Error("should have an empty result")
	}

	// Timeout should be after 1 second
	elapsedCheck(t, start, 1)
}

// TestFutureCancel
func TestFutureCancel(t *testing.T) {
	start := time.Now()

	// Get a context with an explicit cancel function
	ctx, cancel := context.WithCancel(context.Background())

	// Wait a second, and then cancel the future.
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	future := RunAsync(ctx, slowFunction)

	// We should time out with a "context deadline exceeded" error
	res, err := future.Result()
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Error("received unexpected error: ", err)
	}

	// Result should be empty
	if res != "" {
		t.Error("should have an empty result")
	}

	// Timeout should be after 1 second
	elapsedCheck(t, start, 1)
}

func elapsedCheck(t *testing.T, start time.Time, seconds int) {
	elapsed := int(time.Since(start).Seconds())

	if seconds != elapsed {
		t.Errorf("expected %d seconds; got %d\n", seconds, elapsed)
	}
}
