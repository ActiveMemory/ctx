//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"testing"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
)

func TestOpenAIFactoryDefaults(t *testing.T) {
	assertNamedFactory(
		t,
		openAIFactory,
		cfgBackend.NameOpenAI,
		cfgBackend.DefaultEndpointOpenAI,
		cfgBackend.DefaultAPIKeyEnvOpenAI,
	)
}

func TestAnthropicFactoryDefaults(t *testing.T) {
	assertNamedFactory(
		t,
		anthropicFactory,
		cfgBackend.NameAnthropic,
		cfgBackend.DefaultEndpointAnthropic,
		cfgBackend.DefaultAPIKeyEnvAnthropic,
	)
}

func TestOllamaFactoryDefaults(t *testing.T) {
	assertNamedFactory(
		t,
		ollamaFactory,
		cfgBackend.NameOllama,
		cfgBackend.DefaultEndpointOllama,
		"",
	)
}

func TestLMStudioFactoryDefaults(t *testing.T) {
	assertNamedFactory(
		t,
		lmStudioFactory,
		cfgBackend.NameLMStudio,
		cfgBackend.DefaultEndpointLMStudio,
		"",
	)
}

func TestVLLMFactoryDefaults(t *testing.T) {
	assertNamedFactory(
		t,
		vllmFactory,
		cfgBackend.NameVLLM,
		cfgBackend.DefaultEndpointVLLM,
		"",
	)
}

func TestEndpointInfoReportsFactoryDefault(t *testing.T) {
	backend, factoryErr := vllmFactory(Config{})
	if factoryErr != nil {
		t.Fatalf("vllmFactory() error = %v", factoryErr)
	}
	if got := EndpointInfo(backend); got != cfgBackend.DefaultEndpointVLLM {
		t.Fatalf("EndpointInfo() = %q, want %q", got, cfgBackend.DefaultEndpointVLLM)
	}
}

func TestNamedFactoryKeepsConfiguredValues(t *testing.T) {
	backend, factoryErr := openAIFactory(Config{ //nolint:gosec // G101: test fixture, value is an env var name, not a credential
		Endpoint:  "https://example.invalid",
		APIKeyEnv: "CTX_CUSTOM_KEY",
	})
	if factoryErr != nil {
		t.Fatalf("openAIFactory() error = %v", factoryErr)
	}
	wrapped := assertOpenAICompatible(t, backend)
	if wrapped.config.Endpoint != "https://example.invalid" {
		t.Fatalf("Endpoint = %q", wrapped.config.Endpoint)
	}
	if wrapped.config.APIKeyEnv != "CTX_CUSTOM_KEY" {
		t.Fatalf("APIKeyEnv = %q", wrapped.config.APIKeyEnv)
	}
}

func assertNamedFactory(
	t *testing.T,
	factory Factory,
	name string,
	endpoint string,
	apiKeyEnv string,
) {
	t.Helper()
	backend, factoryErr := factory(Config{})
	if factoryErr != nil {
		t.Fatalf("factory() error = %v", factoryErr)
	}
	if backend.Name() != name {
		t.Fatalf("Name() = %q, want %q", backend.Name(), name)
	}
	wrapped := assertOpenAICompatible(t, backend)
	if wrapped.config.Endpoint != endpoint {
		t.Fatalf("Endpoint = %q, want %q", wrapped.config.Endpoint, endpoint)
	}
	if wrapped.config.APIKeyEnv != apiKeyEnv {
		t.Fatalf("APIKeyEnv = %q, want %q", wrapped.config.APIKeyEnv, apiKeyEnv)
	}
}

func assertOpenAICompatible(t *testing.T, backend Backend) openAICompatible {
	t.Helper()
	wrapped, ok := backend.(openAICompatible)
	if !ok {
		t.Fatalf("backend = %T, want openAICompatible", backend)
	}
	return wrapped
}
