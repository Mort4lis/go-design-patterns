package stability

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func counter() UserFunc {
	m := sync.Mutex{}
	count := 0

	return func(ctx context.Context) (string, error) {
		m.Lock()
		count++
		m.Unlock()

		return fmt.Sprintf("%d", count), nil
	}
}

func TestDebounceFirst(t *testing.T) {
	ctx := context.Background()
	debounce := DebounceFirst(counter(), time.Second)

	res, _ := debounce(ctx)
	if res != "1" {
		t.Errorf("wrong debounce result: got %s, want 1", res)
	}

	time.Sleep(900 * time.Millisecond)

	for i := 0; i < 10; i++ {
		res, _ = debounce(ctx)
	}

	if res != "1" {
		t.Errorf("wrong debounce result: got %s, want 1", res)
	}
}

// TestDebounceFirstDataRace tests for data races.
func TestDebounceFirstDataRace(t *testing.T) {
	ctx := context.Background()
	debounce := DebounceFirst(failAfter(1), time.Second)

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

func TestDebounceLast(t *testing.T) {
	var res string

	ctx := context.Background()
	debounce := DebounceLast(counter(), time.Second)

	for i := 0; i < 10; i++ {
		res, _ = debounce(ctx)
	}

	if res != "" {
		t.Errorf("wrong debounce result: got %s, want empty string", res)
	}

	time.Sleep(1100 * time.Millisecond)

	res, _ = debounce(ctx)
	if res != "1" {
		t.Errorf("wrong debounce result: got %s, want 1", res)
	}
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
