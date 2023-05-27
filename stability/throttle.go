package stability

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrToManyCalls = errors.New("too many calls")

func Throttle(fn UserFunc, maxCalls uint, interval time.Duration) UserFunc {
	var once sync.Once
	var mu sync.Mutex
	tokens := maxCalls

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			go func() {
				ticker := time.NewTicker(interval)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						mu.Lock()
						tokens = maxCalls
						mu.Unlock()
					}
				}
			}()
		})

		mu.Lock()
		defer mu.Unlock()

		if tokens <= 0 {
			return "", ErrToManyCalls
		}

		tokens--

		return fn(ctx)
	}
}
