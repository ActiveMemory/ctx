//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// newOllama constructs the Ollama backend. Thin wrapper
// over [newOpenAICompat] that points at the standard
// Ollama local listener (`http://localhost:11434`) and
// sends no auth by default (Ollama is unauthenticated in
// its default install).
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - *openAICompat: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure.
func newOllama(cfg Config) (*openAICompat, error) {
	if cfg.Name == "" {
		cfg.Name = cfgBackend.NameOllama
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = cfgBackend.DefaultEndpointOllama
	}
	return newOpenAICompat(cfg)
}
