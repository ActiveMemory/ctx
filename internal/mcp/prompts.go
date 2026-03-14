//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
)

// promptDefs defines all available MCP prompts.
var promptDefs = []Prompt{
	{
		Name:        mcp.MCPPromptSessionStart,
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartDesc),
	},
	{
		Name:        mcp.MCPPromptAddDecision,
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptAddDecisionDesc),
		Arguments: []PromptArgument{
			{Name: "content", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionTitle), Required: true},
			{Name: "context", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionCtx), Required: true},
			{Name: "rationale", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionRat), Required: true},
			{Name: "consequences", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionConseq), Required: true},
		},
	},
	{
		Name:        mcp.MCPPromptAddLearning,
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptAddLearningDesc),
		Arguments: []PromptArgument{
			{Name: "content", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningTitle), Required: true},
			{Name: "context", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningCtx), Required: true},
			{Name: "lesson", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningLesson), Required: true},
			{Name: "application", Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningApp), Required: true},
		},
	},
	{
		Name:        mcp.MCPPromptReflect,
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptReflectDesc),
	},
	{
		Name:        mcp.MCPPromptCheckpoint,
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointDesc),
	},
}

// handlePromptsList returns all available MCP prompts.
func (s *Server) handlePromptsList(req Request) *Response {
	return s.ok(req.ID, PromptListResult{Prompts: promptDefs})
}

// handlePromptsGet returns the content of a requested prompt.
func (s *Server) handlePromptsGet(req Request) *Response {
	var params GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, errCodeInvalidArg, assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}

	switch params.Name {
	case mcp.MCPPromptSessionStart:
		return s.promptSessionStart(req.ID)
	case mcp.MCPPromptAddDecision:
		return s.promptAddDecision(req.ID, params.Arguments)
	case mcp.MCPPromptAddLearning:
		return s.promptAddLearning(req.ID, params.Arguments)
	case mcp.MCPPromptReflect:
		return s.promptReflect(req.ID)
	case mcp.MCPPromptCheckpoint:
		return s.promptCheckpoint(req.ID)
	default:
		return s.error(req.ID, errCodeNotFound,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPUnknownPrompt), params.Name))
	}
}

// promptSessionStart loads context and provides a session orientation.
func (s *Server) promptSessionStart(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.error(id, errCodeInternal,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	for _, fileName := range ctxCfg.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}
		fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPPromptSectionFormat),
			fileName, string(f.Content))
	}

	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartFooter))

	return s.ok(id, GetPromptResult{
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartResultD),
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: mcp.MCPContentTypeText, Text: sb.String()},
			},
		},
	})
}

// promptAddDecision formats a decision for recording.
func (s *Server) promptAddDecision(
	id json.RawMessage, args map[string]string,
) *Response {
	content := args["content"]
	ctx := args["context"]
	rationale := args["rationale"]
	consequences := args["consequences"]

	fieldFmt := assets.TextDesc(assets.TextDescKeyMCPPromptAddDecisionFieldFmt)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptAddDecisionHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Decision", content, token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Context", ctx, token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Rationale", rationale, token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Consequences", consequences, token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptAddDecisionFooter))

	return s.ok(id, GetPromptResult{
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptAddDecisionResultD),
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: mcp.MCPContentTypeText, Text: sb.String()},
			},
		},
	})
}

// promptAddLearning formats a learning for recording.
func (s *Server) promptAddLearning(
	id json.RawMessage, args map[string]string,
) *Response {
	content := args["content"]
	ctx := args["context"]
	lesson := args["lesson"]
	application := args["application"]

	fieldFmt := assets.TextDesc(assets.TextDescKeyMCPPromptAddLearningFieldFmt)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptAddLearningHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Learning", content, token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Context", ctx, token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Lesson", lesson, token.NewlineLF)
	fmt.Fprintf(&sb, fieldFmt+"%s", "Application", application, token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptAddLearningFooter))

	return s.ok(id, GetPromptResult{
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptAddLearningResultD),
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: mcp.MCPContentTypeText, Text: sb.String()},
			},
		},
	})
}

// promptReflect reviews the current session for outstanding items.
func (s *Server) promptReflect(id json.RawMessage) *Response {
	return s.ok(id, GetPromptResult{
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptReflectResultD),
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: mcp.MCPContentTypeText, Text: assets.TextDesc(assets.TextDescKeyMCPPromptReflectBody)},
			},
		},
	})
}

// promptCheckpoint summarizes progress and prepares for session end.
func (s *Server) promptCheckpoint(id json.RawMessage) *Response {
	pending := s.session.pendingCount()
	adds := totalAdds(s.session.addsPerformed)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointStatsFormat),
		s.session.toolCalls, adds, pending)

	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointSteps))

	return s.ok(id, GetPromptResult{
		Description: assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointResultD),
		Messages: []PromptMessage{
			{
				Role:    "user",
				Content: ToolContent{Type: mcp.MCPContentTypeText, Text: sb.String()},
			},
		},
	})
}
