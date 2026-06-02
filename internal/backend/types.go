//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Backend speaks an OpenAI-compatible HTTP surface to a
// remote inference server.
//
// Implementations are constructed via a [Factory] and
// addressed by their registered type name (e.g., "vllm").
// All methods are safe for concurrent use; backends hold
// per-request state on the stack, never on the receiver.
type Backend interface {
	// Name returns the registered backend type label,
	// e.g., "vllm" or "openai-compatible".
	Name() string
	// Ping performs a reachability check by issuing
	// `GET /v1/models` against the configured endpoint.
	// Returns nil iff the server responds with HTTP 200
	// and a parseable model list. Wraps the underlying
	// transport error otherwise.
	Ping(ctx context.Context) error
	// Models returns the model IDs the backend currently
	// serves, in the order the upstream `/v1/models`
	// response lists them. Used by `ctx ai ping` to surface
	// the first model after a successful reachability
	// check. Wraps the underlying transport / parse error
	// on failure; returns [errBackend.ErrEmptyModels] when
	// the server lists no models.
	Models(ctx context.Context) ([]string, error)
	// Complete issues a single non-streaming chat
	// completion against `/v1/chat/completions`.
	Complete(ctx context.Context, req Request) (Response, error)
}

// Factory constructs a [Backend] from a [Config]. Each
// per-backend implementation registers exactly one
// Factory with the [Registry].
type Factory func(cfg Config) (Backend, error)

// Config is the per-backend slice of `.ctxrc` settings.
//
// Fields:
//   - Name: backend type label, e.g., "vllm". Must match
//     a registered Factory at [Registry.Resolve] time.
//   - Endpoint: base URL, e.g., "http://localhost:8000".
//     Path segments below `/v1/...` are appended by the
//     backend implementation; this value should not end
//     in a slash.
//   - APIKeyEnv: name of the environment variable that
//     holds the bearer token. Empty means no auth header
//     is sent (the vLLM default).
//   - Timeout: per-request deadline. Zero means use the
//     backend implementation's default.
//   - DefaultModel: model ID to use when a request does
//     not specify one. Empty means require explicit
//     model selection at call time.
type Config struct {
	Name         string
	Endpoint     string
	APIKeyEnv    string
	Timeout      time.Duration
	DefaultModel string
}

// Message is a single OpenAI-shape chat message.
//
// Fields:
//   - Role: the OpenAI chat role. Per-backend role
//     labels and the typed constants for them live in
//     `internal/config/backend/` (introduced alongside
//     the first concrete backend implementation).
//   - Content: the message body. Empty content is
//     permitted only for assistant tool-call messages,
//     which are out of scope for v1.
type Message struct {
	Role    string
	Content string
}

// ResponseFormat selects the structured-output mode for
// a request.
//
// Fields:
//   - Type: "json_object" for free-form JSON, or
//     "json_schema" for schema-constrained output.
//   - JSONSchema: the schema payload, used only when
//     Type is "json_schema". The shape mirrors OpenAI's
//     `response_format.json_schema` object; backends
//     translate to vLLM's `guided_json` where the OpenAI
//     contract is not supported.
type ResponseFormat struct {
	Type       string
	JSONSchema map[string]any
}

// Request is the contract floor for `/v1/chat/completions`.
//
// Fields:
//   - Model: target model ID. If empty, the backend uses
//     its [Config.DefaultModel].
//   - Messages: ordered chat history, oldest first.
//   - ResponseFormat: optional structured-output spec.
//     nil means free-form text.
//   - MaxTokens: upper bound on output tokens. Zero
//     leaves the bound at the model default.
//   - Temperature: sampling temperature in [0.0, 2.0].
//     Negative values are reserved as "unset".
type Request struct {
	Model          string
	Messages       []Message
	ResponseFormat *ResponseFormat
	MaxTokens      int
	Temperature    float64
}

// Response captures the assistant's reply.
//
// Fields:
//   - Model: the model ID the server reported (may
//     differ from [Request.Model] when an alias was
//     resolved).
//   - Content: assistant-message text. For
//     schema-constrained outputs this is the JSON string
//     the model emitted; callers parse it.
//   - Raw: the unmodified response body bytes, preserved
//     so callers verifying schema-constrained output
//     have a stable reference even after Content is
//     parsed.
type Response struct {
	Model   string
	Content string
	Raw     []byte
}

// Registry resolves named backend configurations to
// constructed [Backend] instances. It separates two
// concerns kept distinct in the .ctxrc shape: which
// backend *types* the binary knows how to speak
// (Factories, registered at process init), and which
// of those types this project actually configures
// (Configs, loaded from `.ctxrc`).
//
// Fields:
//   - mu: guards factories, configs, and deflt.
//   - factories: type name to constructor.
//   - configs: type name to per-project Config.
//   - deflt: name of the default-selected backend, or
//     empty for "no explicit default."
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
	configs   map[string]Config
	deflt     string
}

// vllm is the vLLM backend. Embeds *openAICompat for the
// wire work and overrides Ping with cold-start retry on
// ECONNREFUSED (the vLLM listener is not yet bound while
// weights are loading; OS returns refused, not 503).
//
// Fields:
//   - openAICompat: embedded generic backend providing
//     Name/Complete and the base Ping.
//   - coldStartWindow: maximum wall-clock during which
//     Ping retries on refused.
//   - coldStartInterval: sleep between retry attempts.
type vllm struct {
	*openAICompat
	coldStartWindow   time.Duration
	coldStartInterval time.Duration
}

// llamacpp is the llama.cpp (llama-server) backend.
// Embeds *openAICompat for the wire work and overrides
// Ping with cold-start retry on ECONNREFUSED (llama-server
// does not bind the listener until the model is fully
// loaded; the OS returns ECONNREFUSED during that window).
//
// Fields:
//   - openAICompat: embedded generic backend providing
//     Name/Complete and the base Ping.
//   - coldStartWindow: maximum wall-clock during which
//     Ping retries on refused.
//   - coldStartInterval: sleep between retry attempts.
type llamacpp struct {
	*openAICompat
	coldStartWindow   time.Duration
	coldStartInterval time.Duration
}

// openAICompat is a generic OpenAI-compatible HTTP
// backend used directly for `openai-compatible` configs
// and embedded by per-vendor wrappers (vllm, openai, ...)
// so the wire work lives in one place.
//
// Fields:
//   - name: registered backend type label.
//   - base: pre-parsed endpoint URL (for url.JoinPath).
//   - apiKey: resolved Authorization Bearer value or "".
//   - timeout: per-request deadline applied via the
//     embedded httpClient.
//   - defaultModel: used when Request.Model is empty.
//   - httpClient: net/http client honoring the timeout.
type openAICompat struct {
	name         string
	base         *url.URL
	apiKey       string
	timeout      time.Duration
	defaultModel string
	httpClient   *http.Client
}
