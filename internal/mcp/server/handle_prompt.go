//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/prompt"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/entity"
)

// handlePromptsList returns all available MCP prompts.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *Response: prompt list result
func (s *Server) handlePromptsList(req proto.Request) *proto.Response {
	return s.ok(req.ID, proto.PromptListResult{Prompts: entity.PromptDefs})
}

// handlePromptsGet returns the content of a requested prompt.
//
// Parameters:
//   - req: the MCP request containing prompt name and arguments
//
// Returns:
//   - *Response: rendered prompt or error
func (s *Server) handlePromptsGet(req proto.Request) *proto.Response {
	var params proto.GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, proto.ErrCodeInvalidArg,
			assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}

	switch params.Name {
	case prompt.SessionStart:
		return s.promptSessionStart(req.ID)
	case prompt.AddDecision:
		return s.promptAddDecision(req.ID, params.Arguments)
	case prompt.AddLearning:
		return s.promptAddLearning(req.ID, params.Arguments)
	case prompt.Reflect:
		return s.promptReflect(req.ID)
	case prompt.Checkpoint:
		return s.promptCheckpoint(req.ID)
	default:
		return s.error(
			req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPUnknownPrompt),
				params.Name,
			),
		)
	}
}

// promptSessionStart loads context and provides session orientation.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: rendered session start prompt with context files
func (s *Server) promptSessionStart(
	id json.RawMessage,
) *proto.Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.error(id, proto.ErrCodeInternal,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	var sb strings.Builder
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartHeader),
	)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	for _, fileName := range ctxCfg.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}
		_, _ = fmt.Fprintf(
			&sb,
			assets.TextDesc(assets.TextDescKeyMCPPromptSectionFormat),
			fileName, string(f.Content),
		)
	}

	sb.WriteString(token.NewlineLF)
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptSessionStartFooter),
	)

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptSessionStartResultD,
		),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: sb.String(),
				},
			},
		},
	})
}

// buildEntryPrompt renders a structured entry prompt (decision or
// learning) from the given spec and returns the formatted response.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - spec: entry prompt specification (header, footer, fields)
//
// Returns:
//   - *proto.Response: formatted entry prompt
func (s *Server) buildEntryPrompt(
	id json.RawMessage, spec entity.EntryPromptSpec,
) *proto.Response {
	fieldFmt := assets.TextDesc(spec.FieldFmtK)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(spec.HeaderKey))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	for _, f := range spec.Fields {
		_, _ = fmt.Fprintf(
			&sb,
			fieldFmt, assets.TextDesc(f.LabelKey), f.Value,
		)
	}
	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(spec.FooterKey))

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(spec.ResultDKey),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: sb.String(),
				},
			},
		},
	})
}

// promptAddDecision formats a decision for recording.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: prompt arguments (content, context, rationale,
//     consequences)
//
// Returns:
//   - *Response: formatted decision prompt
func (s *Server) promptAddDecision(
	id json.RawMessage, args map[string]string,
) *proto.Response {
	return s.buildEntryPrompt(id, entity.EntryPromptSpec{
		HeaderKey:  assets.TextDescKeyMCPPromptAddDecisionHeader,
		FooterKey:  assets.TextDescKeyMCPPromptAddDecisionFooter,
		FieldFmtK:  assets.TextDescKeyMCPPromptAddDecisionFieldFmt,
		ResultDKey: assets.TextDescKeyMCPPromptAddDecisionResultD,
		Fields: []entity.EntryField{
			{LabelKey: assets.TextDescKeyMCPPromptLabelDecision,
				Value: args[field.Content]},
			{LabelKey: assets.TextDescKeyMCPPromptLabelContext,
				Value: args[cli.AttrContext]},
			{LabelKey: assets.TextDescKeyMCPPromptLabelRationale,
				Value: args[cli.AttrRationale]},
			{LabelKey: assets.TextDescKeyMCPPromptLabelConsequences,
				Value: args[cli.AttrConsequences]},
		},
	})
}

// promptAddLearning formats a learning for recording.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: prompt arguments (content, context, lesson,
//     application)
//
// Returns:
//   - *Response: formatted learning prompt
func (s *Server) promptAddLearning(
	id json.RawMessage, args map[string]string,
) *proto.Response {
	return s.buildEntryPrompt(id, entity.EntryPromptSpec{
		HeaderKey:  assets.TextDescKeyMCPPromptAddLearningHeader,
		FooterKey:  assets.TextDescKeyMCPPromptAddLearningFooter,
		FieldFmtK:  assets.TextDescKeyMCPPromptAddLearningFieldFmt,
		ResultDKey: assets.TextDescKeyMCPPromptAddLearningResultD,
		Fields: []entity.EntryField{
			{LabelKey: assets.TextDescKeyMCPPromptLabelLearning,
				Value: args[field.Content]},
			{LabelKey: assets.TextDescKeyMCPPromptLabelContext,
				Value: args[cli.AttrContext]},
			{LabelKey: assets.TextDescKeyMCPPromptLabelLesson,
				Value: args[cli.AttrLesson]},
			{LabelKey: assets.TextDescKeyMCPPromptLabelApplication,
				Value: args[cli.AttrApplication]},
		},
	})
}

// promptReflect reviews the current session for outstanding items.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: reflection prompt text
func (s *Server) promptReflect(id json.RawMessage) *proto.Response {
	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptReflectResultD),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: assets.TextDesc(
						assets.TextDescKeyMCPPromptReflectBody,
					),
				},
			},
		},
	})
}

// promptCheckpoint summarizes progress and prepares for session
// end.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: checkpoint prompt with session stats
func (s *Server) promptCheckpoint(
	id json.RawMessage,
) *proto.Response {
	pending := s.session.PendingCount()
	adds := totalAdds(s.session.AddsPerformed)

	var sb strings.Builder
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointHeader),
	)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointStatsFormat),
		s.session.ToolCalls, adds, pending,
	)

	sb.WriteString(token.NewlineLF)
	sb.WriteString(
		assets.TextDesc(assets.TextDescKeyMCPPromptCheckpointSteps),
	)

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointResultD,
		),
		Messages: []proto.PromptMessage{
			{
				Role: prompt.RoleUser,
				Content: proto.ToolContent{
					Type: mime.ContentTypeText,
					Text: sb.String(),
				},
			},
		},
	})
}
