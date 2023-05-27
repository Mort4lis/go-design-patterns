package stability

import (
	"context"
	"log"
	"time"
)

func Retry(fn Circuit, maxRetries int, delay time.Duration) Circuit {
	return func(ctx context.Context) (string, error) {
		for attempt := 1; ; attempt++ {
			result, err := fn(ctx)
			if err == nil || attempt >= maxRetries {
				return result, err
			}

			log.Printf("Attempt %d failed, retry after %v", attempt, delay)

			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}

				return "", ctx.Err()
			case <-timer.C:
			}
		}
	}
}
