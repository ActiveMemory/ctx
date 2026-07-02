//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import "time"

// Backend names for built-in backend factories.
const (
	NameAnthropic        = "anthropic"
	NameLMStudio         = "lmstudio"
	NameOllama           = "ollama"
	NameOpenAI           = "openai"
	NameOpenAICompatible = "openai-compatible"
	NameVLLM             = "vllm"
)

// Default endpoints for built-in backend factories.
const (
	DefaultEndpointAnthropic = "https://api.anthropic.com"
	DefaultEndpointLMStudio  = "http://localhost:1234"
	DefaultEndpointOllama    = "http://localhost:11434"
	DefaultEndpointOpenAI    = "https://api.openai.com"
	DefaultEndpointVLLM      = "http://localhost:8000"
)

// Default API key environment variables for backend factories.
const (
	DefaultAPIKeyEnvAnthropic = "ANTHROPIC_API_KEY"
	DefaultAPIKeyEnvOpenAI    = "OPENAI_API_KEY"
)

// HTTP constants for OpenAI-compatible backends.
const (
	AuthorizationBearerPrefix = "Bearer "
	ChatCompletionsPath       = "/v1/chat/completions"
	ContentTypeJSON           = "application/json"
	MaxResponseBytes          = 1 << 20
	ResponseFormatJSONSchema  = "json_schema"
	HeaderAuthorization       = "Authorization"
	HeaderContentType         = "Content-Type"
	HTTPMethodGet             = "GET"
	HTTPMethodPost            = "POST"
	ModelsPath                = "/v1/models"
	RoleUser                  = "user"
)

// DefaultTimeout is the fallback timeout for backend HTTP calls.
const DefaultTimeout = 30 * time.Second

// Error messages for backend registry and HTTP failures.
const (
	ErrBadRequest            = "backend request failed: "
	ErrDuplicateRegistration = "backend already registered: "
	ErrFactory               = "backend factory failed: "
	ErrInvalidEndpoint       = "backend endpoint invalid: "
	ErrInvalidResponseShape  = "invalid structured response shape"
	ErrMissingBackend        = "backend not registered: "
	ErrMultipleBackends      = "multiple backends configured; pass --backend " +
		"or set backends.default"
	ErrNoBackendConfigured = "no backend configured; run ctx setup --backend"
	FmtUpstreamStatusBody  = "%d: %s"
	ErrUnreachable         = "backend unreachable: "
	ErrUpstream            = "backend upstream returned: "
)
