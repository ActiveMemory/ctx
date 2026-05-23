//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// newLMStudio constructs the LM Studio backend. Thin
// wrapper over [newOpenAICompat] that points at the
// standard LM Studio local listener
// (`http://localhost:1234`) and sends no auth by default
// (LM Studio is unauthenticated in its default install).
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - *openAICompat: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure.
func newLMStudio(cfg Config) (*openAICompat, error) {
	if cfg.Name == "" {
		cfg.Name = cfgBackend.NameLMStudio
	}
	if cfg.Endpoint == "" {
		cfg.Endpoint = cfgBackend.DefaultEndpointLMStudio
	}
	return newOpenAICompat(cfg)
}
