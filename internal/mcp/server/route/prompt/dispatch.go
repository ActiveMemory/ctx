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
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	defPrompt "github.com/ActiveMemory/ctx/internal/mcp/server/def/prompt"
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
		req.ID, proto.PromptListResult{Prompts: defPrompt.Defs},
	)
}

// DispatchGet unmarshals prompt params and dispatches to the
// appropriate prompt builder.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory and session
//   - req: the MCP request containing prompt name and arguments
//
// Returns:
//   - *proto.Response: rendered prompt or error
func DispatchGet(
	d *entity.MCPDeps, req proto.Request,
) *proto.Response {
	var params proto.GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return out.ErrResponse(req.ID, cfgSchema.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams))
	}

	switch params.Name {
	case prompt.SessionStart:
		return sessionStart(req.ID, d.ContextDir)
	case prompt.AddDecision:
		return addDecision(req.ID, params.Arguments)
	case prompt.AddLearning:
		return addLearning(req.ID, params.Arguments)
	case prompt.Reflect:
		return reflect(req.ID)
	case prompt.Checkpoint:
		return checkpoint(
			req.ID,
			d.Session.ToolCalls,
			d.Session.AddsPerformed,
			d.Session.PendingCount(),
		)
	default:
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeNotFound,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPErrUnknownPrompt),
				params.Name,
			),
		)
	}
}
