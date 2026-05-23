//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// newAnthropic constructs the Anthropic backend. Thin
// wrapper over [newOpenAICompat] that points at
// api.anthropic.com by default and reads ANTHROPIC_API_KEY
// from the environment. Anthropic ships an
// OpenAI-compatible `/v1/chat/completions` endpoint
// alongside the native Messages API, so the contract
// floor works as-is; the wrapper does not need its own
// wire shape. The native Messages API path is out of
// scope for this thin-wrapper layer.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - *openAICompat: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure.
func newAnthropic(cfg Config) (*openAICompat, error) {
	if cfg.Name == "" {
		cfg.Name = cfgBackend.NameAnthropic
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = cfgBackend.DefaultEndpointAnthropic
	}
	if cfg.APIKeyEnv == "" {
		cfg.APIKeyEnv = cfgBackend.EnvAnthropicAPIKey
	}
	return newOpenAICompat(cfg)
}
