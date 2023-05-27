package stability

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func emulateTransientError() Circuit {
	var count int

	return func(ctx context.Context) (string, error) {
		count++

		if count <= 3 {
			return "", errors.New("error")
		} else {
			return fmt.Sprintf("%d", count), nil
		}
	}
}

func TestRetry(t *testing.T) {
	ctx := context.Background()
	r := Retry(emulateTransientError(), 5, 50*time.Millisecond)

	res, err := r(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res != "4" {
		t.Errorf("wrong result: got %s, want 3", res)
	}
}
