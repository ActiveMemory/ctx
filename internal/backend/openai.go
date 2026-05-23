//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// newOpenAI constructs the OpenAI backend. Thin wrapper
// over [newOpenAICompat] that sets vendor-specific
// defaults (name, endpoint, API-key env var) when the
// caller's [Config] leaves them empty. The user's .ctxrc
// values always win — defaults only fill gaps.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - *openAICompat: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure.
func newOpenAI(cfg Config) (*openAICompat, error) {
	if cfg.Name == "" {
		cfg.Name = cfgBackend.NameOpenAI
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = cfgBackend.DefaultEndpointOpenAI
	}
	if cfg.APIKeyEnv == "" {
		cfg.APIKeyEnv = cfgBackend.EnvOpenAIAPIKey
	}
	return newOpenAICompat(cfg)
}
