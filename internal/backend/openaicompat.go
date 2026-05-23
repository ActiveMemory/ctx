//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	cfgBackend "github.com/ActiveMemory/ctx/internal/config/backend"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
)

// Name implements [Backend].
//
// Returns:
//   - string: the registered backend type label.
func (b *openAICompat) Name() string { return b.name }

// Ping implements [Backend] by issuing
// `GET <endpoint>/v1/models` and asserting HTTP 200.
// Any transport-level failure surfaces as
// [errBackend.ErrUnreachable]; any non-200 surfaces as
// [errBackend.ErrUnhealthyStatus] with a truncated body
// excerpt.
//
// Parameters:
//   - ctx: caller-provided context for cancellation.
//
// Returns:
//   - error: nil on success, typed sentinel on failure.
func (b *openAICompat) Ping(ctx context.Context) error {
	u := b.base.JoinPath(cfgBackend.PathV1, cfgBackend.PathModels).String()
	req, reqErr := http.NewRequestWithContext(
		ctx, http.MethodGet, u, nil,
	)
	if reqErr != nil {
		return errBackend.Unreachable(b.name, u, reqErr)
	}
	if b.apiKey != "" {
		req.Header.Set(
			cfgBackend.HeaderAuthorization,
			cfgBackend.AuthBearerPrefix+b.apiKey,
		)
	}
	resp, doErr := b.httpClient.Do(req)
	if doErr != nil {
		return errBackend.Unreachable(b.name, u, doErr)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(
			io.LimitReader(resp.Body, cfgBackend.BodyExcerptLimit),
		)
		return errBackend.UnhealthyStatus(
			b.name, resp.StatusCode, string(body),
		)
	}
	return nil
}

// Models implements [Backend] by issuing
// `GET <endpoint>/v1/models` and parsing the response
// `data[].id` array. Reachability errors flow through
// [errBackend.ErrUnreachable]; non-200 responses through
// [errBackend.ErrUnhealthyStatus]; a 200 response with an
// empty `data` array surfaces as
// [errBackend.ErrEmptyModels] so callers can distinguish
// "reachable but unusable" from a true transport failure.
//
// Parameters:
//   - ctx: caller-provided context for cancellation.
//
// Returns:
//   - []string: ordered list of model IDs as the server
//     reported them.
//   - error: nil on success, typed sentinel on failure.
func (b *openAICompat) Models(ctx context.Context) ([]string, error) {
	u := b.base.JoinPath(cfgBackend.PathV1, cfgBackend.PathModels).String()
	req, reqErr := http.NewRequestWithContext(
		ctx, http.MethodGet, u, nil,
	)
	if reqErr != nil {
		return nil, errBackend.Unreachable(b.name, u, reqErr)
	}
	if b.apiKey != "" {
		req.Header.Set(
			cfgBackend.HeaderAuthorization,
			cfgBackend.AuthBearerPrefix+b.apiKey,
		)
	}
	resp, doErr := b.httpClient.Do(req)
	if doErr != nil {
		return nil, errBackend.Unreachable(b.name, u, doErr)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, readErr := io.ReadAll(
		io.LimitReader(resp.Body, cfgBackend.MaxResponseBytes),
	)
	if readErr != nil {
		return nil, errBackend.ReadResponse(b.name, readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errBackend.UnhealthyStatus(
			b.name, resp.StatusCode, excerpt(raw),
		)
	}
	var parsed modelsResponse
	if unmarshalErr := json.Unmarshal(raw, &parsed); unmarshalErr != nil {
		return nil, errBackend.ParseResponse(b.name, unmarshalErr)
	}
	if len(parsed.Data) == 0 {
		return nil, errBackend.EmptyModels(b.name)
	}
	ids := make([]string, len(parsed.Data))
	for i, m := range parsed.Data {
		ids[i] = m.ID
	}
	return ids, nil
}

// Complete implements [Backend] by issuing
// `POST <endpoint>/v1/chat/completions`. Maps the public
// [Request] to the OpenAI wire shape, validates the
// response, and returns the first choice's message
// content along with the raw body bytes.
//
// Parameters:
//   - ctx: caller-provided context for cancellation.
//   - r: the chat-completion request.
//
// Returns:
//   - Response: populated on success.
//   - error: typed sentinel on failure.
func (b *openAICompat) Complete(
	ctx context.Context, r Request,
) (Response, error) {
	model := r.Model
	if model == "" {
		model = b.defaultModel
	}
	if model == "" {
		return Response{}, errBackend.MissingModel(b.name)
	}
	payload := newChatRequest(model, r)
	body, mErr := json.Marshal(payload)
	if mErr != nil {
		return Response{}, errBackend.MarshalRequest(b.name, mErr)
	}
	u := b.base.JoinPath(
		cfgBackend.PathV1, cfgBackend.PathChat, cfgBackend.PathCompletions,
	).String()
	req, reqErr := http.NewRequestWithContext(
		ctx, http.MethodPost, u, bytes.NewReader(body),
	)
	if reqErr != nil {
		return Response{}, errBackend.Unreachable(b.name, u, reqErr)
	}
	req.Header.Set(cfgBackend.HeaderContentType, cfgBackend.MIMEJSON)
	if b.apiKey != "" {
		req.Header.Set(
			cfgBackend.HeaderAuthorization,
			cfgBackend.AuthBearerPrefix+b.apiKey,
		)
	}
	resp, doErr := b.httpClient.Do(req)
	if doErr != nil {
		return Response{}, errBackend.Unreachable(b.name, u, doErr)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, readErr := io.ReadAll(
		io.LimitReader(resp.Body, cfgBackend.MaxResponseBytes),
	)
	if readErr != nil {
		return Response{}, errBackend.ReadResponse(b.name, readErr)
	}
	if resp.StatusCode != http.StatusOK {
		return Response{}, errBackend.UpstreamStatus(
			b.name, resp.StatusCode, excerpt(raw),
		)
	}
	var parsed chatResponse
	if unmarshalErr := json.Unmarshal(raw, &parsed); unmarshalErr != nil {
		return Response{}, errBackend.ParseResponse(b.name, unmarshalErr)
	}
	if len(parsed.Choices) == 0 {
		return Response{}, errBackend.EmptyChoices(b.name)
	}
	return Response{
		Model:   parsed.Model,
		Content: parsed.Choices[0].Message.Content,
		Raw:     raw,
	}, nil
}
