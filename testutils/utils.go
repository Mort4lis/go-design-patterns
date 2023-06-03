package testutils

import (
	"runtime"
	"testing"
)

func DetectGoroutineLeeks(t *testing.T, fn func()) {
	t.Helper()

	beforeNum := runtime.NumGoroutine()
	fn()
	afterNum := runtime.NumGoroutine()

	if beforeNum < afterNum {
		buf := make([]byte, 4096)
		runtime.Stack(buf, true)

		t.Log(string(buf))
		t.Errorf(
			"gorotines leek is detected: before running test was %d, now %d",
			beforeNum, afterNum,
		)
	}
}
