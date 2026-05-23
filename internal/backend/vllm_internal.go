//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"errors"
	"syscall"
	"time"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// newVLLM constructs a vLLM backend. The wire work is
// delegated to the embedded openAICompat; vllm adds
// cold-start retry on ECONNREFUSED for Ping. The vLLM
// docs describe a multi-minute weight-load window during
// which the OpenAI-compatible listener is not yet bound;
// during that window the OS returns ECONNREFUSED (not
// HTTP 503), so HTTP-status-based retry is wrong.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - *vllm: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure (missing/invalid endpoint).
func newVLLM(cfg Config) (*vllm, error) {
	if cfg.Name == "" {
		cfg.Name = cfgBackend.NameVLLM
	}
	inner, err := newOpenAICompat(cfg)
	if err != nil {
		return nil, err
	}
	return &vllm{
		openAICompat:      inner,
		coldStartWindow:   cfgBackend.DefaultColdStartWindow,
		coldStartInterval: cfgBackend.DefaultColdStartInterval,
	}, nil
}

// coldStartRetry calls ping repeatedly while the error is
// a dial-refused (the listener is not yet bound) and the
// window has not elapsed. Any successful call, non-dial
// error, window expiry, or context cancellation returns
// immediately.
//
// Parameters:
//   - ctx: caller-provided context.
//   - ping: the inner reachability check.
//   - window: maximum wall-clock to keep retrying.
//   - interval: sleep between attempts.
//
// Returns:
//   - error: nil on success, last error on failure.
func coldStartRetry(
	ctx context.Context,
	ping func(context.Context) error,
	window, interval time.Duration,
) error {
	deadline := time.Now().Add(window)
	for {
		err := ping(ctx)
		if err == nil {
			return nil
		}
		if !isDialRefused(err) {
			return err
		}
		if time.Now().After(deadline) {
			return err
		}
		select {
		case <-time.After(interval):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// isDialRefused reports whether err's chain contains
// syscall.ECONNREFUSED. Used to distinguish "server not
// yet listening" (retry) from "server returned non-200"
// (don't retry).
//
// Parameters:
//   - err: error returned by openAICompat.Ping.
//
// Returns:
//   - bool: true when the underlying cause is ECONNREFUSED.
func isDialRefused(err error) bool {
	return errors.Is(err, syscall.ECONNREFUSED)
}
