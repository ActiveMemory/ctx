//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"

// builtinFactory returns a built-in backend factory by name.
//
// Parameters:
//   - name: backend implementation name
//
// Returns:
//   - Factory: backend factory
//   - bool: whether the factory exists
func builtinFactory(name string) (Factory, bool) {
	switch name {
	case cfgBackend.NameOpenAICompatible:
		return openAICompatibleFactory, true
	case cfgBackend.NameVLLM:
		return vllmFactory, true
	case cfgBackend.NameOpenAI:
		return openAIFactory, true
	case cfgBackend.NameAnthropic:
		return anthropicFactory, true
	case cfgBackend.NameOllama:
		return ollamaFactory, true
	case cfgBackend.NameLMStudio:
		return lmStudioFactory, true
	default:
		return nil, false
	}
}
