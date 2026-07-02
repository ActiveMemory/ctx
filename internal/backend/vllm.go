//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"

// vllmFactory creates the canonical local vLLM backend.
//
// Parameters:
//   - config: backend configuration
//
// Returns:
//   - Backend: configured vLLM backend
//   - error: always nil; request validation happens during calls
func vllmFactory(config Config) (Backend, error) {
	return openAICompatible{
		name:   cfgBackend.NameVLLM,
		config: defaultConfig(config, cfgBackend.DefaultEndpointVLLM, ""),
	}, nil
}
