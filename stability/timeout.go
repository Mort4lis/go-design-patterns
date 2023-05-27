package stability

import "context"

type SlowFunc func(s string) (string, error)

type WithContext func(ctx context.Context, s string) (string, error)

func Timeout(fn SlowFunc) WithContext {
	return func(ctx context.Context, s string) (string, error) {
		resCh := make(chan string)
		errCh := make(chan error)

		go func() {
			res, err := fn(s)
			resCh <- res
			errCh <- err
		}()

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case res := <-resCh:
			return res, <-errCh
		}
	}
}
