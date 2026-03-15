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
)

// promptDefs defines all available MCP prompts.
var promptDefs = []proto.Prompt{
	{
		Name: prompt.SessionStart,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptSessionStartDesc),
	},
	{
		Name: prompt.AddDecision,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptAddDecisionDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrRationale,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionRat),
				Required:    true,
			},
			{
				Name:        cli.AttrConsequences,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgDecisionConseq),
				Required:    true,
			},
		},
	},
	{
		Name: prompt.AddLearning,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptAddLearningDesc),
		Arguments: []proto.PromptArgument{
			{
				Name:        field.Content,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningTitle),
				Required:    true,
			},
			{
				Name:        cli.AttrContext,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningCtx),
				Required:    true,
			},
			{
				Name:        cli.AttrLesson,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningLesson),
				Required:    true,
			},
			{
				Name:        cli.AttrApplication,
				Description: assets.TextDesc(assets.TextDescKeyMCPPromptArgLearningApp),
				Required:    true,
			},
		},
	},
	{
		Name: prompt.Reflect,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptReflectDesc),
	},
	{
		Name: prompt.Checkpoint,
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointDesc),
	},
}

// handlePromptsList returns all available MCP prompts.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *Response: prompt list result
func (s *Server) handlePromptsList(req proto.Request) *proto.Response {
	return s.ok(req.ID, proto.PromptListResult{Prompts: promptDefs})
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
		return s.error(req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPUnknownPrompt),
				params.Name))
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
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptSessionStartHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	for _, fileName := range ctxCfg.ReadOrder {
		f := ctx.File(fileName)
		if f == nil || f.IsEmpty {
			continue
		}
		_, _ = fmt.Fprintf(&sb,
			assets.TextDesc(assets.TextDescKeyMCPPromptSectionFormat),
			fileName, string(f.Content))
	}

	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptSessionStartFooter))

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptSessionStartResultD),
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
	content := args[field.Content]
	ctx := args[cli.AttrContext]
	rationale := args[cli.AttrRationale]
	consequences := args[cli.AttrConsequences]

	fieldFmt := assets.TextDesc(
		assets.TextDescKeyMCPPromptAddDecisionFieldFmt)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptAddDecisionHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(assets.TextDescKeyMCPPromptLabelDecision),
		content)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(assets.TextDescKeyMCPPromptLabelContext),
		ctx)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(assets.TextDescKeyMCPPromptLabelRationale),
		rationale)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(
			assets.TextDescKeyMCPPromptLabelConsequences),
		consequences)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptAddDecisionFooter))

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptAddDecisionResultD),
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
	content := args[field.Content]
	ctx := args[cli.AttrContext]
	lesson := args[cli.AttrLesson]
	application := args[cli.AttrApplication]

	fieldFmt := assets.TextDesc(
		assets.TextDescKeyMCPPromptAddLearningFieldFmt)

	var sb strings.Builder
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptAddLearningHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(assets.TextDescKeyMCPPromptLabelLearning),
		content)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(assets.TextDescKeyMCPPromptLabelContext),
		ctx)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(assets.TextDescKeyMCPPromptLabelLesson),
		lesson)
	_, _ = fmt.Fprintf(&sb, fieldFmt,
		assets.TextDesc(
			assets.TextDescKeyMCPPromptLabelApplication),
		application)
	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptAddLearningFooter))

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptAddLearningResultD),
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
						assets.TextDescKeyMCPPromptReflectBody),
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
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptCheckpointHeader))
	sb.WriteString(token.NewlineLF)
	sb.WriteString(token.NewlineLF)

	_, _ = fmt.Fprintf(&sb,
		assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointStatsFormat),
		s.session.ToolCalls, adds, pending)

	sb.WriteString(token.NewlineLF)
	sb.WriteString(assets.TextDesc(
		assets.TextDescKeyMCPPromptCheckpointSteps))

	return s.ok(id, proto.GetPromptResult{
		Description: assets.TextDesc(
			assets.TextDescKeyMCPPromptCheckpointResultD),
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
