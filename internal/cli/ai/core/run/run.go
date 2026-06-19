//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package run

import (
	"context"
	"encoding/json"

	backendPkg "github.com/ActiveMemory/ctx/internal/backend"
	cfgAI "github.com/ActiveMemory/ctx/internal/config/ai"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errBackend "github.com/ActiveMemory/ctx/internal/err/backend"
	ctxio "github.com/ActiveMemory/ctx/internal/io"
)

// Ping resolves and pings a configured backend.
//
// Parameters:
//   - ctx: request context
//   - backendName: optional backend selector
//
// Returns:
//   - PingResult: backend, endpoint, and first model
//   - error: backend resolution or ping failure
func Ping(ctx context.Context, backendName string) (PingResult, error) {
	resolved, resolveErr := resolve(backendName)
	if resolveErr != nil {
		return PingResult{}, resolveErr
	}
	info, pingErr := backendPkg.PingInfo(ctx, resolved.backend)
	if pingErr != nil {
		return PingResult{}, pingErr
	}
	return PingResult{
		Backend:    resolved.name,
		Endpoint:   backendPkg.EndpointInfo(resolved.backend),
		FirstModel: info.FirstModel,
	}, nil
}

// Propose generates and writes a validation-only proposal artifact.
//
// Parameters:
//   - ctx: request context
//   - input: input file path
//   - backendName: optional backend selector
//   - emit: comma-separated proposal kinds
//
// Returns:
//   - string: artifact path
//   - error: backend, completion, validation, or write failure
func Propose(
	ctx context.Context,
	input string,
	backendName string,
	emit string,
) (string, error) {
	resolved, resolveErr := resolve(backendName)
	if resolveErr != nil {
		return "", resolveErr
	}
	data, readErr := ctxio.SafeReadUserFile(input)
	if readErr != nil {
		return "", readErr
	}
	kinds := splitEmit(emit)
	response, completeErr := resolved.backend.Complete(
		ctx,
		backendPkg.Request{
			Prompt: cfgAI.PromptPrefix + emit + token.NewlineLF + string(data),
			Schema: cfgAI.SchemaMinimal,
		},
	)
	if completeErr != nil {
		return "", completeErr
	}
	decoded := map[string]any{}
	decodeErr := json.Unmarshal([]byte(response.Text), &decoded)
	if decodeErr != nil {
		return "", errBackend.BadRequest{
			Name:  resolved.name,
			Cause: decodeErr,
		}
	}
	artifact := ProposalArtifact{
		Kind:     cfgAI.KindProposedPatch,
		Backend:  resolved.name,
		Model:    response.Model,
		Input:    input,
		Emit:     kinds,
		Status:   cfgAI.StatusProposed,
		Response: decoded,
	}
	return writeArtifact(artifact)
}
