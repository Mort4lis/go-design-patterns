package stability

import (
	"context"
	"sync"
	"time"
)

func DebounceFirst(fn UserFunc, d time.Duration) UserFunc {
	var (
		mu              sync.Mutex
		result          string
		err             error
		thresholdCallAt time.Time
	)

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer func() {
			thresholdCallAt = time.Now().Add(d)
			mu.Unlock()
		}()

		if time.Now().Before(thresholdCallAt) {
			return result, err
		}

		result, err = fn(ctx)
		return result, err
	}
}

func DebounceLast(fn UserFunc, d time.Duration) UserFunc {
	var (
		mu              sync.Mutex
		once            sync.Once
		result          string
		err             error
		thresholdCallAt time.Time
	)

	return func(ctx context.Context) (string, error) {
		mu.Lock()
		defer mu.Unlock()

		thresholdCallAt = time.Now().Add(d)

		once.Do(func() {
			go func() {
				ticker := time.NewTicker(100 * time.Millisecond)

				defer func() {
					ticker.Stop()

					mu.Lock()
					once = sync.Once{}
					mu.Unlock()
				}()

				for {
					select {
					case <-ctx.Done():
						mu.Lock()
						result, err = "", ctx.Err()
						mu.Unlock()

						return
					case <-ticker.C:
						mu.Lock()
						if time.Now().Before(thresholdCallAt) {
							mu.Unlock()
							continue
						}

						result, err = fn(ctx)
						mu.Unlock()

						return
					}
				}
			}()
		})

		return result, err
	}
}
