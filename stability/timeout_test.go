package stability

import (
	"context"
	"errors"
	"testing"
	"time"
)

func createWork(d time.Duration) SlowFunc {
	return func(s string) (string, error) {
		time.Sleep(d)
		return s, nil
	}
}

func TestTimeout(t *testing.T) {
	workFunc := createWork(100 * time.Millisecond)
	timeout := Timeout(workFunc)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	res, err := timeout(ctx, "")
	t.Logf("result=%s, err=%v", res, err)

	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected timeout error, got %v", err)
	}
}

func TestTimeoutWithoutFired(t *testing.T) {
	workFunc := createWork(50 * time.Millisecond)
	timeout := Timeout(workFunc)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	res, err := timeout(ctx, "")
	t.Logf("result=%s, err=%v", res, err)

	if errors.Is(err, context.DeadlineExceeded) {
		t.Error("unexpected timeout error")
	}
}
