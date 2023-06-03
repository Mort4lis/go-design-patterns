package testutils

import (
	"testing"
	"time"
)

type ChanConfig struct {
	timer   *time.Timer
	timeout *time.Duration
}

type ChanOption func(cfg *ChanConfig)

func WithTimer(t *time.Timer) ChanOption {
	return func(cfg *ChanConfig) {
		cfg.timer = t
	}
}

func WithDuration(d time.Duration) ChanOption {
	return func(cfg *ChanConfig) {
		cfg.timeout = &d
	}
}

func ReadChan(t *testing.T, ch <-chan any, opts ...ChanOption) any {
	t.Helper()

	var cfg ChanConfig
	timeoutChan := timeoutChanFromOpts(&cfg, opts...)

	select {
	case val := <-ch:
		return val
	case <-timeoutChan:
		t.Fatalf("timeout exceeded while reading the channel")
	}

	return nil
}

func WriteChan(t *testing.T, ch chan<- any, val any, opts ...ChanOption) {
	t.Helper()

	var cfg ChanConfig
	timeoutChan := timeoutChanFromOpts(&cfg, opts...)

	select {
	case ch <- val:
	case <-timeoutChan:
		t.Fatalf("timeout exceeded while writing the channel")
	}
}

func CheckClosedChan(t *testing.T, ch <-chan any, opts ...ChanOption) {
	t.Helper()

	var cfg ChanConfig
	timeoutChan := timeoutChanFromOpts(&cfg, opts...)

	isOutClosed := false

	select {
	case _, ok := <-ch:
		isOutClosed = !ok
	case <-timeoutChan:
	}

	if !isOutClosed {
		t.Errorf("channel isn't closed, but expected")
	}
}

func timeoutChanFromOpts(cfg *ChanConfig, opts ...ChanOption) <-chan time.Time {
	for _, opt := range opts {
		opt(cfg)
	}

	switch {
	case cfg.timer != nil:
		return cfg.timer.C
	case cfg.timeout != nil:
		return time.After(*cfg.timeout)
	}

	return nil
}
