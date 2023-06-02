package concurrency

import (
	"context"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var (
		got  int
		want int
	)
	genCh := Generator(ctx)

	for i := 0; i < 5; i++ {
		want += i
		got += <-genCh
	}

	cancel()

	if got != want {
		t.Errorf("wrong result: got %d, want %d", got, want)
	}

	isGenClosed := false

	select {
	case _, ok := <-genCh:
		if !ok {
			isGenClosed = true
		}
	case <-time.After(50 * time.Millisecond):
	}

	if !isGenClosed {
		t.Errorf("generator isn't closed when context done")
	}
}
