package stability

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold int) Circuit {
	var mu sync.RWMutex
	failures := 0
	lastAttemptAt := time.Now()

	return func(ctx context.Context) (string, error) {
		mu.RLock()

		d := failures - failureThreshold
		if d >= 0 {
			shouldRetryAt := lastAttemptAt.Add(2 * time.Second << d)
			if shouldRetryAt.After(time.Now()) {
				mu.RUnlock()
				return "", errors.New("service is unavailable")
			}
		}

		mu.RUnlock()
		resp, err := circuit(ctx)

		mu.Lock()
		defer mu.Unlock()

		lastAttemptAt = time.Now()

		if err != nil {
			failures++
			return resp, err
		}

		failures = 0

		return resp, nil
	}
}
