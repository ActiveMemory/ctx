//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	prompt2 "github.com/ActiveMemory/ctx/internal/mcp/server/prompt"
)

// handlePromptsList returns all available MCP prompts.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *proto.Response: prompt list result
func (s *Server) handlePromptsList(req proto.Request) *proto.Response {
	return out.OkResponse(req.ID, proto.PromptListResult{Prompts: entity.PromptDefs})
}

// handlePromptsGet returns the content of a requested prompt.
//
// Parameters:
//   - req: the MCP request containing prompt name and arguments
//
// Returns:
//   - *proto.Response: rendered prompt or error
func (s *Server) handlePromptsGet(req proto.Request) *proto.Response {
	var params proto.GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return out.ErrResponse(req.ID, proto.ErrCodeInvalidArg,
			assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}

	switch params.Name {
	case prompt.SessionStart:
		return prompt2.SessionStart(req.ID, s.handler.ContextDir)
	case prompt.AddDecision:
		return prompt2.AddDecision(req.ID, params.Arguments)
	case prompt.AddLearning:
		return prompt2.AddLearning(req.ID, params.Arguments)
	case prompt.Reflect:
		return prompt2.Reflect(req.ID)
	case prompt.Checkpoint:
		return prompt2.Checkpoint(
			req.ID,
			s.handler.Session.ToolCalls,
			s.handler.Session.AddsPerformed,
			s.handler.Session.PendingCount(),
		)
	default:
		return out.ErrResponse(
			req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPUnknownPrompt),
				params.Name,
			),
		)
	}
}
