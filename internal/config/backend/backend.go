//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import "time"

// Registered backend type labels. These are the `name:`
// values accepted in the `.ctxrc backends:` table and the
// strings the registry uses as map keys.
const (
	// NameOpenAICompat is the registered backend type
	// label for the generic OpenAI-compatible backend.
	NameOpenAICompat = "openai-compatible"
	// NameVLLM is the registered backend type label for
	// the vLLM wrapper (adds cold-start retry).
	NameVLLM = "vllm"
	// NameOpenAI is the registered backend type label for
	// the OpenAI wrapper.
	NameOpenAI = "openai"
	// NameAnthropic is the registered backend type label
	// for the Anthropic wrapper. Anthropic ships an
	// OpenAI-compatible `/v1/chat/completions` endpoint;
	// the wrapper inherits the contract floor and points
	// at api.anthropic.com by default.
	NameAnthropic = "anthropic"
	// NameOllama is the registered backend type label for
	// the Ollama wrapper.
	NameOllama = "ollama"
	// NameLMStudio is the registered backend type label
	// for the LM Studio wrapper.
	NameLMStudio = "lmstudio"
)

// Default endpoints for each per-vendor wrapper. Applied
// only when Config.Endpoint is empty so the user override
// in `.ctxrc` always wins.
const (
	// DefaultEndpointOpenAI is the OpenAI public API
	// base URL.
	DefaultEndpointOpenAI = "https://api.openai.com"
	// DefaultEndpointAnthropic is the Anthropic public
	// API base URL.
	DefaultEndpointAnthropic = "https://api.anthropic.com"
	// DefaultEndpointOllama is the Ollama default local
	// listener address.
	DefaultEndpointOllama = "http://localhost:11434"
	// DefaultEndpointLMStudio is the LM Studio default
	// local listener address.
	DefaultEndpointLMStudio = "http://localhost:1234"
	// DefaultEndpointVLLM is the vLLM canonical default
	// local listener address. Used by `ctx setup --backend
	// vllm` when --endpoint is not passed; vLLM
	// deployments routinely vary the port so the user can
	// always override.
	DefaultEndpointVLLM = "http://localhost:8000"
)

// Default API-key environment variable names for each
// per-vendor wrapper. Applied only when Config.APIKeyEnv
// is empty so the user override in `.ctxrc` always wins.
// Ollama and LM Studio default to no auth — empty values
// mean no Authorization header is sent.
const (
	// EnvOpenAIAPIKey is the canonical OpenAI key env var.
	EnvOpenAIAPIKey = "OPENAI_API_KEY"
	// EnvAnthropicAPIKey is the canonical Anthropic key
	// env var.
	EnvAnthropicAPIKey = "ANTHROPIC_API_KEY"
)

// URL scheme tokens accepted in Config.Endpoint. The
// constants prevent accidental literal duplication in the
// validator and document the scheme allowlist explicitly.
const (
	// SchemeHTTP is the only insecure scheme accepted in
	// backend Endpoint URLs. https is preferred but http
	// stays valid for localhost vLLM endpoints.
	SchemeHTTP = "http"
	// SchemeHTTPS is the secure scheme accepted in
	// backend Endpoint URLs.
	SchemeHTTPS = "https"
	// EndpointTrimChar is the trailing character stripped
	// from Endpoint before url.Parse so JoinPath produces
	// clean URLs across `http://x:8000` and
	// `http://x:8000/` inputs.
	EndpointTrimChar = "/"
)

// HTTP path segments for the OpenAI-compatible surface
// that the contract floor uses. Joined via url.URL.JoinPath
// so the resulting URLs respect the configured Endpoint
// base path (any future user-mounted prefix is preserved).
const (
	// PathV1 is the version segment shared by /v1/models
	// and /v1/chat/completions.
	PathV1 = "v1"
	// PathModels is the path segment for the models
	// listing endpoint used by Ping.
	PathModels = "models"
	// PathChat is the path segment for the chat namespace.
	PathChat = "chat"
	// PathCompletions is the path segment for the
	// chat-completions endpoint.
	PathCompletions = "completions"
)

// HTTP header names, value prefixes, MIME types, and
// excerpt formatting tokens. Kept here so the wire layer
// in `internal/backend/` reads no string literals.
const (
	// HeaderAuthorization is the HTTP header name used to
	// send the bearer token to the backend.
	HeaderAuthorization = "Authorization"
	// HeaderContentType is the HTTP header name used to
	// declare JSON request bodies.
	HeaderContentType = "Content-Type"
	// AuthBearerPrefix is the prefix used in the
	// Authorization header value.
	AuthBearerPrefix = "Bearer "
	// MIMEJSON is the Content-Type value used for outgoing
	// chat-completion request bodies.
	MIMEJSON = "application/json"
	// BodyTruncSuffix is appended to error-body excerpts
	// that exceed [BodyExcerptLimit].
	BodyTruncSuffix = "..."
)

// Timeout defaults applied when the per-backend Config
// leaves the value zero. Long enough to cover real
// cold-start windows on local inference servers, short
// enough that a hung backend does not pin the agent.
const (
	// DefaultRequestTimeout is the per-request deadline
	// used when Config.Timeout is zero.
	DefaultRequestTimeout = 30 * time.Second
	// DefaultColdStartWindow is the maximum wall-clock
	// window during which vllm.Ping retries on
	// ECONNREFUSED.
	DefaultColdStartWindow = 90 * time.Second
	// DefaultColdStartInterval is how long Ping sleeps
	// between retry attempts during the cold-start window.
	DefaultColdStartInterval = 500 * time.Millisecond
)

// Response-body size limits applied by the wire layer.
// Caps protect against pathologically large upstream
// responses and keep error messages readable.
const (
	// MaxResponseBytes caps how much of an upstream
	// response body the chat-completions path will read.
	// Prevents a pathologically large response from
	// exhausting Go's heap.
	MaxResponseBytes = 8 * 1024 * 1024
	// BodyExcerptLimit caps how much of an error response
	// is copied into user-facing error text. 512 bytes is
	// enough to surface model-not-found / rate-limit /
	// auth error bodies without dumping multi-KB HTML
	// error pages from misconfigured proxies.
	BodyExcerptLimit = 512
)
