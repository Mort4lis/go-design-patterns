package concurrency

import (
	"context"
	"sync"
)

type Future interface {
	Result() (string, error)
}

type InnerFuture struct {
	once     sync.Once
	wg       sync.WaitGroup
	result   string
	err      error
	resultCh <-chan string
	errCh    <-chan error
}

func (f *InnerFuture) Result() (string, error) {
	f.once.Do(func() {
		f.wg.Add(1)
		defer f.wg.Done()

		f.result, f.err = <-f.resultCh, <-f.errCh
	})

	f.wg.Wait()

	return f.result, f.err
}

type SlowFunc func() (string, error)

func RunAsync(ctx context.Context, fn SlowFunc) Future {
	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	go func() {
		result, err := fn()
		resultCh <- result
		errCh <- err
	}()

	go func() {
		<-ctx.Done()

		resultCh <- ""
		errCh <- ctx.Err()
	}()

	return &InnerFuture{resultCh: resultCh, errCh: errCh}
}
