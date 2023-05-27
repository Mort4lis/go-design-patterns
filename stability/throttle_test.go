package stability

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// callsCountFunction returns a function that increments the int that
// callCounter points to whenever it is run.
func callsCountFunction(callCounter *int) UserFunc {
	return func(ctx context.Context) (string, error) {
		*callCounter++
		return fmt.Sprintf("call %d", *callCounter), nil
	}
}

// TestThrottleMax1 tests whether a max of 1 call per duration is respected.
func TestThrottleMax1(t *testing.T) {
	const max uint = 1

	callsCounter := 0
	fn := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(fn, max, time.Second)

	for i := 0; i < 100; i++ {
		_, _ = throttle(ctx)
	}

	if callsCounter == 0 {
		t.Error("test is broken; got", callsCounter)
	}

	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

// TestThrottleMax10 tests whether a max of 10 calls per duration is respected.
func TestThrottleMax10(t *testing.T) {
	const max uint = 10

	callsCounter := 0
	fn := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(fn, max, time.Second)

	for i := 0; i < 100; i++ {
		_, _ = throttle(ctx)
	}

	if callsCounter == 0 {
		t.Error("test is broken; got", callsCounter)
	}

	if callsCounter > int(max) {
		t.Error("max is broken; got", callsCounter)
	}
}

// TestThrottleCallFrequency5Seconds tests whether a Throttle with a max of 1
// and a duration of 1 second called every 250ms for 5 seconds will be called
// exactly 5 times.
func TestThrottleCallFrequency5Seconds(t *testing.T) {
	callsCounter := 0
	fn := callsCountFunction(&callsCounter)

	ctx := context.Background()
	throttle := Throttle(fn, 1, time.Second)

	// make a call every 1/4 second for 5 seconds.
	tickCounts := 0
	ticker := time.NewTicker(250 * time.Millisecond).C

	for range ticker {
		tickCounts++

		s, e := throttle(ctx)
		if e != nil {
			t.Log("Error:", e)
		} else {
			t.Log("output:", s)
		}

		if tickCounts >= 20 {
			break
		}
	}

	if callsCounter != 5 {
		t.Error("expected 5; got", callsCounter)
	}
}

// TestThrottleContextTimeout tests whether a Throttle will return an error
// when its context is canceled.
func TestThrottleContextTimeout(t *testing.T) {
	callsCounter := 0
	fn := callsCountFunction(&callsCounter)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	throttle := Throttle(fn, 1, time.Second)

	s, e := throttle(ctx)
	if e != nil {
		t.Error("unexpected error:", e)
	} else {
		t.Log("output:", s)
	}

	// Wait for timeout
	time.Sleep(300 * time.Millisecond)

	_, e = throttle(ctx)
	if e != nil {
		t.Log("got expected error:", e)
	} else {
		t.Error("didn't get expected error")
	}
}
