//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	promptdef "github.com/ActiveMemory/ctx/internal/mcp/server/def/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// DispatchList returns all available prompts.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *proto.Response: prompt list response
func DispatchList(req proto.Request) *proto.Response {
	return out.OkResponse(
		req.ID, proto.PromptListResult{Prompts: promptdef.Defs},
	)
}

// DispatchGet unmarshals prompt params and dispatches to the
// appropriate prompt builder.
//
// Parameters:
//   - h: handler for context directory and session state
//   - req: the MCP request containing prompt name and arguments
//
// Returns:
//   - *proto.Response: rendered prompt or error
func DispatchGet(
	h *handler.Handler, req proto.Request,
) *proto.Response {
	var params proto.GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return out.ErrResponse(req.ID, proto.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams))
	}

	switch params.Name {
	case prompt.SessionStart:
		return sessionStart(req.ID, h.ContextDir)
	case prompt.AddDecision:
		return addDecision(req.ID, params.Arguments)
	case prompt.AddLearning:
		return addLearning(req.ID, params.Arguments)
	case prompt.Reflect:
		return reflect(req.ID)
	case prompt.Checkpoint:
		return checkpoint(
			req.ID,
			h.Session.ToolCalls,
			h.Session.AddsPerformed,
			h.Session.PendingCount(),
		)
	default:
		return out.ErrResponse(
			req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPErrUnknownPrompt),
				params.Name,
			),
		)
	}
}
