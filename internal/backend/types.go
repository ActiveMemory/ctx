//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"encoding/json"
)

// Backend is the AI completion backend contract.
type Backend interface {
	Name() string
	Ping(ctx context.Context) error
	Complete(ctx context.Context, req Request) (Response, error)
}

// Factory builds a backend from registry configuration.
type Factory func(Config) (Backend, error)

// Config carries rc-backed backend settings into backend factories.
//
// Fields:
//   - Name: configured backend name
//   - Type: registered backend implementation type
//   - Endpoint: backend HTTP endpoint, when applicable
//   - APIKeyEnv: environment variable name for credentials
//   - Timeout: request timeout duration string from configuration
//   - DefaultModel: model selected when a request omits one
type Config struct {
	Name         string
	Type         string
	Endpoint     string
	APIKeyEnv    string
	Timeout      string
	DefaultModel string
}

// Request describes a backend completion request.
//
// Fields:
//   - Model: model override for this request
//   - Prompt: input prompt text
//   - Schema: optional JSON schema for structured output
type Request struct {
	Model  string
	Prompt string
	Schema Schema
}

// Schema carries a structured-output schema definition.
type Schema struct {
	Name   string
	Schema json.RawMessage
}

// Response describes a backend completion response.
//
// Fields:
//   - Model: model that produced the response
//   - Text: completion text
//   - Raw: optional raw provider payload
type Response struct {
	Model      string
	Text       string
	Raw        []byte
	FirstModel string
}

// Registry stores backend factories and resolves configured backends.
type Registry struct {
	factories   map[string]Factory
	configs     map[string]Config
	defaultName string
}

// openAICompatible is an OpenAI-compatible HTTP backend.
type openAICompatible struct {
	name   string
	config Config
}

// modelsResponse is the OpenAI-compatible /v1/models response.
type modelsResponse struct {
	Data []modelInfo `json:"data"`
}

// modelInfo is a single model entry.
type modelInfo struct {
	ID string `json:"id"`
}

// chatRequest is the OpenAI-compatible chat completion request.
type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

// responseFormat encodes an OpenAI-compatible structured output schema.
type responseFormat struct {
	Type       string                `json:"type"`
	JSONSchema *responseFormatSchema `json:"json_schema,omitempty"`
}

// responseFormatSchema is the structured output schema wrapper.
type responseFormatSchema struct {
	Name   string          `json:"name"`
	Strict bool            `json:"strict"`
	Schema json.RawMessage `json:"schema"`
}

// chatMessage is a chat message in an OpenAI-compatible payload.
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse is the OpenAI-compatible chat completion response.
type chatResponse struct {
	Model   string       `json:"model"`
	Choices []chatChoice `json:"choices"`
}

// chatChoice is one chat completion candidate.
type chatChoice struct {
	Message chatMessage `json:"message"`
}
