//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// newLlamaCpp constructs a llama.cpp backend. The wire
// work is delegated to the embedded openAICompat;
// llamacpp adds cold-start retry on ECONNREFUSED for
// Ping. llama-server does not bind its HTTP listener
// until the model weights are fully loaded; during that
// window the OS returns ECONNREFUSED, so the same
// cold-start retry logic used by vLLM applies here.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - *llamacpp: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure (missing/invalid endpoint).
func newLlamaCpp(cfg Config) (*llamacpp, error) {
	if cfg.Name == "" {
		cfg.Name = cfgBackend.NameLlamaCpp
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = cfgBackend.DefaultEndpointLlamaCpp
	}
	inner, err := newOpenAICompat(cfg)
	if err != nil {
		return nil, err
	}
	return &llamacpp{
		openAICompat:      inner,
		coldStartWindow:   cfgBackend.DefaultColdStartWindow,
		coldStartInterval: cfgBackend.DefaultColdStartInterval,
	}, nil
}

// Ping implements [Backend] with llama.cpp-specific
// cold-start retry. Delegates to the embedded
// openAICompat.Ping via [coldStartRetry] which retries
// while the failure is a dial-refused error (server has
// not yet bound the listener) and the cold-start window
// has not elapsed.
//
// Parameters:
//   - ctx: caller-provided context for cancellation.
//
// Returns:
//   - error: nil on success, typed sentinel on failure.
func (b *llamacpp) Ping(ctx context.Context) error {
	return coldStartRetry(
		ctx,
		b.openAICompat.Ping,
		b.coldStartWindow,
		b.coldStartInterval,
	)
}
