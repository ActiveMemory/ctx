//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

// vllmFactory adapts [newVLLM] (which returns the
// concrete `*vllm` type) to the [Factory] signature
// (which returns the [Backend] interface). Adapter
// pattern needed because Go does not auto-convert a
// `func(Config) (*vllm, error)` to a
// `func(Config) (Backend, error)`.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - Backend: the constructed vLLM backend on success.
//   - error: typed err/backend sentinel on failure.
func vllmFactory(cfg Config) (Backend, error) {
	return newVLLM(cfg)
}

// openAICompatFactory adapts [newOpenAICompat] to the
// [Factory] signature.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - Backend: the constructed backend on success.
//   - error: typed err/backend sentinel on failure.
func openAICompatFactory(cfg Config) (Backend, error) {
	return newOpenAICompat(cfg)
}

// openAIFactory adapts [newOpenAI] to the [Factory]
// signature.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - Backend: the constructed backend on success.
//   - error: typed err/backend sentinel on failure.
func openAIFactory(cfg Config) (Backend, error) {
	return newOpenAI(cfg)
}

// anthropicFactory adapts [newAnthropic] to the
// [Factory] signature.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - Backend: the constructed backend on success.
//   - error: typed err/backend sentinel on failure.
func anthropicFactory(cfg Config) (Backend, error) {
	return newAnthropic(cfg)
}

// ollamaFactory adapts [newOllama] to the [Factory]
// signature.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - Backend: the constructed backend on success.
//   - error: typed err/backend sentinel on failure.
func ollamaFactory(cfg Config) (Backend, error) {
	return newOllama(cfg)
}

// lmStudioFactory adapts [newLMStudio] to the [Factory]
// signature.
//
// Parameters:
//   - cfg: per-project backend config.
//
// Returns:
//   - Backend: the constructed backend on success.
//   - error: typed err/backend sentinel on failure.
func lmStudioFactory(cfg Config) (Backend, error) {
	return newLMStudio(cfg)
}
