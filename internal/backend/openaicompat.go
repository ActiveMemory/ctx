//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"context"
	"encoding/json"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

// Name returns the backend implementation name.
//
// Returns:
//   - string: backend name
func (backend openAICompatible) Name() string {
	return backend.name
}

// Ping checks model-list reachability.
//
// Parameters:
//   - ctx: request context
//
// Returns:
//   - error: unreachable, upstream, or response decode failure
func (backend openAICompatible) Ping(ctx context.Context) error {
	_, pingErr := backend.models(ctx)
	return pingErr
}

// Complete requests a chat completion.
//
// Parameters:
//   - ctx: request context
//   - req: completion request
//
// Returns:
//   - Response: completion response
//   - error: unreachable, upstream, encode, or decode failure
func (backend openAICompatible) Complete(
	ctx context.Context,
	req Request,
) (Response, error) {
	model := req.Model
	if model == "" {
		model = backend.config.DefaultModel
	}
	payload := chatRequest{
		Model: model,
		Messages: []chatMessage{{
			Role:    cfgBackend.RoleUser,
			Content: req.Prompt,
		}},
	}
	if req.Schema.Name != "" || req.Schema.Schema != nil {
		payload.ResponseFormat = &responseFormat{
			Type: cfgBackend.ResponseFormatJSONSchema,
			JSONSchema: &responseFormatSchema{
				Name:   req.Schema.Name,
				Strict: true,
				Schema: req.Schema.Schema,
			},
		}
	}
	body, marshalErr := json.Marshal(payload)
	if marshalErr != nil {
		return Response{}, errBackend.BadRequest{
			Name:  backend.name,
			Cause: marshalErr,
		}
	}
	raw, doErr := backend.do(
		ctx,
		cfgBackend.HTTPMethodPost,
		cfgBackend.ChatCompletionsPath,
		body,
	)
	if doErr != nil {
		return Response{}, doErr
	}
	var decoded chatResponse
	if decodeErr := json.Unmarshal(raw, &decoded); decodeErr != nil {
		return Response{}, errBackend.BadRequest{
			Name:  backend.name,
			Cause: decodeErr,
		}
	}
	if decoded.Model == "" ||
		len(decoded.Choices) == 0 ||
		decoded.Choices[0].Message.Content == "" {
		return Response{}, errBackend.BadRequest{
			Name:  backend.name,
			Cause: errBackend.InvalidResponseShape(),
		}
	}
	text := ""
	if len(decoded.Choices) > 0 {
		text = decoded.Choices[0].Message.Content
	}
	return Response{Model: decoded.Model, Text: text, Raw: raw}, nil
}
