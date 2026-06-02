//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

// RegisterAll registers all seven built-in backend
// factories with the supplied Registry. Called by the
// CLI before resolving any backend; keeps the per-vendor
// factory functions unexported (so they cannot be reached
// directly by external callers) while still letting the
// CLI compose a fully-loaded Registry.
//
// Parameters:
//   - r: the Registry to populate. Typically a fresh one
//     from [New].
//
// Returns:
//   - error: any Register failure (duplicate type name);
//     under normal use the registered names are distinct
//     by construction so this returns nil.
func RegisterAll(r *Registry) error {
	pairs := []struct {
		name    string
		factory Factory
	}{
		{cfgBackend.NameVLLM, vllmFactory},
		{cfgBackend.NameOpenAICompat, openAICompatFactory},
		{cfgBackend.NameOpenAI, openAIFactory},
		{cfgBackend.NameAnthropic, anthropicFactory},
		{cfgBackend.NameOllama, ollamaFactory},
		{cfgBackend.NameLMStudio, lmStudioFactory},
		{cfgBackend.NameLlamaCpp, llamacppFactory},
	}
	for _, p := range pairs {
		if err := r.Register(p.name, p.factory); err != nil {
			return err
		}
	}
	return nil
}
