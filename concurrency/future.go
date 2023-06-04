package concurrency

import (
	"context"
	"sync"
)

// Future pattern provides a placeholder for a value that will be generated
// by an asynchronous process. Channels can be often used in a similar way, for example:
//
//	 func ConcurrentInverse() <-chan any {
//		 result := make(chan any)
//		 go func() {
//			 // blocking operation
//			 result <- ""
//		 }()
//		 return result
//	 }
//
// but this is not desirable for something like a public API, because caller can call
// several ConcurrentInverse with the incorrect way, for example:
//
//	func Example() any, any {
//	  return <-ConcurrentInverse(), <-ConcurrentInverse()
//	}
//
// as a result, calls of ConcurrentInverse will be executed serially, requiring twice the runtime.
// Future pattern encapsulates this complexity in an API that provides the consumer with a simple interface
// whose method can be called normally.
type Future interface {
	// Result returns the result blocking the call until the result is ready.
	Result() (string, error)
}

// InnerFuture implements Future.
// It retrieves ready values from the channels, caches and returns them. If the values aren’t available
// on the channels, the request blocks. If they have already been retrieved, the cached values are returned.
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

// SlowFunc is some kind of blocking function that should be performed asynchronously.
type SlowFunc func() (string, error)

// RunAsync is a wrapper function that asynchronously performs SlowFunc and returns Future.
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
