//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

// newOpenAICompat is the constructor for the generic
// openai-compatible backend. Returns the concrete type so
// package peers (vllm, future per-vendor wrappers) can
// embed it. A public Factory wrapper will land in the CLI
// registration step (Task 7) once a non-test consumer
// exists.
//
// Parameters:
//   - cfg: the per-project backend config.
//
// Returns:
//   - *openAICompat: concrete backend.
//   - error: typed err/backend sentinel on validation
//     failure.
func newOpenAICompat(cfg Config) (*openAICompat, error) {
	name := cfg.Name
	if name == "" {
		name = cfgBackend.NameOpenAICompat
	}
	if cfg.Endpoint == "" {
		return nil, errBackend.MissingEndpoint(name)
	}
	u, parseErr := url.Parse(cfg.Endpoint)
	if parseErr != nil ||
		(u.Scheme != cfgBackend.SchemeHTTP && u.Scheme != cfgBackend.SchemeHTTPS) ||
		u.Host == "" {
		return nil, errBackend.InvalidEndpoint(name, cfg.Endpoint)
	}
	apiKey := ""
	if cfg.APIKeyEnv != "" {
		apiKey = os.Getenv(cfg.APIKeyEnv)
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = cfgBackend.DefaultRequestTimeout
	}
	trimmed := strings.TrimRight(cfg.Endpoint, cfgBackend.EndpointTrimChar)
	base, _ := url.Parse(trimmed)
	return &openAICompat{
		name:         name,
		base:         base,
		apiKey:       apiKey,
		timeout:      timeout,
		defaultModel: cfg.DefaultModel,
		httpClient:   &http.Client{Timeout: timeout},
	}, nil
}

// excerpt returns a UTF-8 safe head of a response body
// for use in error messages, with a configurable
// truncation suffix appended when the body exceeds
// [cfgBackend.BodyExcerptLimit].
//
// Parameters:
//   - body: the response body bytes.
//
// Returns:
//   - string: a truncated excerpt.
func excerpt(body []byte) string {
	if len(body) <= cfgBackend.BodyExcerptLimit {
		return string(body)
	}
	return string(body[:cfgBackend.BodyExcerptLimit]) + cfgBackend.BodyTruncSuffix
}

// newChatRequest maps the public [Request] to the OpenAI
// wire shape. The public Temperature uses the
// "negative = unset" convention; this helper translates
// that to a nullable pointer in the wire format so the
// caller can omit the field instead of always sending 0.
//
// Parameters:
//   - model: the resolved model id.
//   - r: the public request.
//
// Returns:
//   - chatRequest: the wire payload, ready for
//     json.Marshal.
func newChatRequest(model string, r Request) chatRequest {
	msgs := make([]chatRequestMessage, len(r.Messages))
	for i, m := range r.Messages {
		msgs[i] = chatRequestMessage(m)
	}
	out := chatRequest{
		Model:     model,
		Messages:  msgs,
		MaxTokens: r.MaxTokens,
	}
	if r.Temperature >= 0 {
		t := r.Temperature
		out.Temperature = &t
	}
	if r.ResponseFormat != nil {
		out.ResponseFormat = &chatRequestResponseFormat{
			Type:       r.ResponseFormat.Type,
			JSONSchema: r.ResponseFormat.JSONSchema,
		}
	}
	return out
}
