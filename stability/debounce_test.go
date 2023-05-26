package stability

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func counter() Circuit {
	m := sync.Mutex{}
	count := 0

	return func(ctx context.Context) (string, error) {
		m.Lock()
		count++
		m.Unlock()

		return fmt.Sprintf("%d", count), nil
	}
}

// TestDebounceFirstDataRace tests for data races.
func TestDebounceFirstDataRace(t *testing.T) {
	ctx := context.Background()

	circuit := failAfter(1)
	debounce := DebounceFirst(circuit, time.Second)

	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	time.Sleep(time.Second * 2)

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := debounce(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	wg.Wait()
}

// TestDebounceLastDataRace tests for data races.
func TestDebounceLastDataRace(t *testing.T) {
	ctx := context.Background()
	debounce := DebounceLast(counter(), time.Second)
	wg := sync.WaitGroup{}

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			res, err := debounce(ctx)
			t.Logf("attempt %d: result=%s, err=%v", count, res, err)
		}(count)
	}

	wg.Wait()

	t.Log("Waiting 2 seconds")

	time.Sleep(time.Second * 2)

	for count := 1; count <= 10; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			res, err := debounce(ctx)
			t.Logf("attempt %d: result=%s, err=%v", count, res, err)
		}(count)
	}

	wg.Wait()
}
