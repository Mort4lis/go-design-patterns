package stability

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

// failAfter returns a function matching the UserFunc type that returns an
// error after it's been called more than threshold times.
func failAfter(threshold int) UserFunc {
	count := 0

	// Service function. Fails after 5 tries.
	return func(ctx context.Context) (string, error) {
		count++

		if count > threshold {
			return "", errors.New("INTENTIONAL FAIL!")
		}

		return "Success", nil
	}
}

func waitAndContinue() UserFunc {
	return func(ctx context.Context) (string, error) {
		time.Sleep(time.Second)

		if rand.Int()%2 == 0 {
			return "success", nil
		}

		return "Failed", fmt.Errorf("forced failure")
	}
}

// TestBreaker tests that the CircuitBreaker function automatically closes and reopens.
func TestCircuitBreaker(t *testing.T) {
	ctx := context.Background()
	// A circuit breaker that opens after one failed attempt.
	cb := CircuitBreaker(failAfter(5), 1)

	circuitOpen := false
	doesCircuitOpen := false
	doesCircuitReclose := false
	count := 0

	for range time.NewTicker(time.Second).C {
		_, err := cb(ctx)

		if err != nil {
			// Does the circuit open?
			if strings.HasPrefix(err.Error(), "service is unavailable") {
				if !circuitOpen {
					circuitOpen = true
					doesCircuitOpen = true

					t.Log("circuit has opened")
				}
			} else {
				// Does it close again?
				if circuitOpen {
					circuitOpen = false
					doesCircuitReclose = true

					t.Log("circuit has automatically closed")
				}
			}
		} else {
			t.Log("circuit closed and operational")
		}

		count++
		if count >= 10 {
			break
		}
	}

	if !doesCircuitOpen {
		t.Error("circuit didn't appear to open")
	}

	if !doesCircuitReclose {
		t.Error("circuit didn't appear to close after time")
	}
}

// TestCircuitBreakerDataRace tests for data races.
func TestCircuitBreakerDataRace(t *testing.T) {
	ctx := context.Background()
	cb := CircuitBreaker(waitAndContinue(), 1)

	wg := sync.WaitGroup{}

	for count := 1; count <= 20; count++ {
		wg.Add(1)

		go func(count int) {
			defer wg.Done()

			time.Sleep(50 * time.Millisecond)

			_, err := cb(ctx)

			t.Logf("attempt %d: err=%v", count, err)
		}(count)
	}

	wg.Wait()
}
