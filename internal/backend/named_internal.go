//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"

// openAIFactory creates the named OpenAI backend wrapper.
//
// Parameters:
//   - config: backend configuration
//
// Returns:
//   - Backend: configured OpenAI-compatible backend
//   - error: always nil; request validation happens during calls
func openAIFactory(config Config) (Backend, error) {
	return openAICompatible{
		name: cfgBackend.NameOpenAI,
		config: defaultConfig(
			config,
			cfgBackend.DefaultEndpointOpenAI,
			cfgBackend.DefaultAPIKeyEnvOpenAI,
		),
	}, nil
}

// anthropicFactory creates the named Anthropic backend wrapper.
//
// Anthropic currently uses the OpenAI-compatible floor. Messages API
// specialization is deferred until backend capability detection exists.
//
// Parameters:
//   - config: backend configuration
//
// Returns:
//   - Backend: configured OpenAI-compatible backend
//   - error: always nil; request validation happens during calls
func anthropicFactory(config Config) (Backend, error) {
	return openAICompatible{
		name: cfgBackend.NameAnthropic,
		config: defaultConfig(
			config,
			cfgBackend.DefaultEndpointAnthropic,
			cfgBackend.DefaultAPIKeyEnvAnthropic,
		),
	}, nil
}

// ollamaFactory creates the named Ollama backend wrapper.
//
// Parameters:
//   - config: backend configuration
//
// Returns:
//   - Backend: configured OpenAI-compatible backend
//   - error: always nil; request validation happens during calls
func ollamaFactory(config Config) (Backend, error) {
	return openAICompatible{
		name: cfgBackend.NameOllama,
		config: defaultConfig(
			config,
			cfgBackend.DefaultEndpointOllama,
			"",
		),
	}, nil
}

// lmStudioFactory creates the named LM Studio backend wrapper.
//
// Parameters:
//   - config: backend configuration
//
// Returns:
//   - Backend: configured OpenAI-compatible backend
//   - error: always nil; request validation happens during calls
func lmStudioFactory(config Config) (Backend, error) {
	return openAICompatible{
		name: cfgBackend.NameLMStudio,
		config: defaultConfig(
			config,
			cfgBackend.DefaultEndpointLMStudio,
			"",
		),
	}, nil
}

// defaultConfig applies endpoint and auth defaults.
//
// Parameters:
//   - config: caller-supplied config
//   - endpoint: default endpoint when config omits one
//   - apiKeyEnv: default API key env var when config omits one
//
// Returns:
//   - Config: config with defaults applied
func defaultConfig(config Config, endpoint string, apiKeyEnv string) Config {
	if config.Endpoint == "" {
		config.Endpoint = endpoint
	}
	if config.APIKeyEnv == "" {
		config.APIKeyEnv = apiKeyEnv
	}
	return config
}
