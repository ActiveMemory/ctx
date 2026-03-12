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
	"github.com/ActiveMemory/ctx/internal/cli/complete"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entry"
)

// toolDefs defines all available MCP tools.
var toolDefs = []Tool{
	{
		Name:        config.MCPToolStatus,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolStatusDesc),
		InputSchema: InputSchema{Type: config.SchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name:        config.MCPToolAdd,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolAddDesc),
		InputSchema: InputSchema{
			Type: config.SchemaObject,
			Properties: map[string]Property{
				config.AttrType: {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropType),
					Enum:        []string{"task", "decision", "learning", "convention"},
				},
				"content": {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropContent),
				},
				"priority": {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropPriority),
					Enum:        []string{"high", "medium", "low"},
				},
				config.AttrContext: {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropContext),
				},
				config.AttrRationale: {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropRationale),
				},
				config.AttrConsequences: {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropConseq),
				},
				config.AttrLesson: {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropLesson),
				},
				config.AttrApplication: {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropApplication),
				},
			},
			Required: []string{config.AttrType, "content"},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name:        config.MCPToolComplete,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolCompleteDesc),
		InputSchema: InputSchema{
			Type: config.SchemaObject,
			Properties: map[string]Property{
				"query": {
					Type:        config.SchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{"query"},
		},
		Annotations: &ToolAnnotations{IdempotentHint: true},
	},
	{
		Name:        config.MCPToolDrift,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolDriftDesc),
		InputSchema: InputSchema{Type: config.SchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
}

// handleToolsList returns all available MCP tools.
func (s *Server) handleToolsList(req Request) *Response {
	return s.ok(req.ID, ToolListResult{Tools: toolDefs})
}

// handleToolsCall dispatches a tool call to the appropriate handler.
func (s *Server) handleToolsCall(req Request) *Response {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(req.ID, errCodeInvalidArg, assets.TextDesc(assets.TextDescKeyMCPInvalidParams))
	}

	switch params.Name {
	case config.MCPToolStatus:
		return s.toolStatus(req.ID)
	case config.MCPToolAdd:
		return s.toolAdd(req.ID, params.Arguments)
	case config.MCPToolComplete:
		return s.toolComplete(req.ID, params.Arguments)
	case config.MCPToolDrift:
		return s.toolDrift(req.ID)
	default:
		return s.error(req.ID, errCodeNotFound,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPUnknownTool), params.Name))
	}
}

// toolStatus loads context and returns a status summary.
func (s *Server) toolStatus(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPStatusContextFormat), ctx.Dir)
	_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPStatusFilesFormat), len(ctx.Files))
	_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPStatusTokensFormat), ctx.TotalTokens)

	for _, f := range ctx.Files {
		status := assets.TextDesc(assets.TextDescKeyMCPStatusOK)
		if f.IsEmpty {
			status = assets.TextDesc(assets.TextDescKeyMCPStatusEmpty)
		}
		_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPStatusFileFormat),
			f.Name, f.Tokens, status)
	}

	return s.toolOK(id, sb.String())
}

// toolAdd adds an entry to a context file.
func (s *Server) toolAdd(
	id json.RawMessage, args map[string]interface{},
) *Response {
	entryType, _ := args[config.AttrType].(string)
	content, _ := args["content"].(string)

	if entryType == "" || content == "" {
		return s.toolError(id, assets.TextDesc(assets.TextDescKeyMCPTypeContentRequired))
	}

	params := entry.Params{
		Type:       entryType,
		Content:    content,
		ContextDir: s.contextDir,
	}

	// Optional fields.
	if v, ok := args["priority"].(string); ok {
		params.Priority = v
	}
	if v, ok := args["context"].(string); ok {
		params.Context = v
	}
	if v, ok := args["rationale"].(string); ok {
		params.Rationale = v
	}
	if v, ok := args["consequences"].(string); ok {
		params.Consequences = v
	}
	if v, ok := args["lesson"].(string); ok {
		params.Lesson = v
	}
	if v, ok := args["application"].(string); ok {
		params.Application = v
	}

	// Validate required fields.
	if vErr := entry.Validate(params, nil); vErr != nil {
		return s.toolError(id, vErr.Error())
	}

	if wErr := entry.Write(params); wErr != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPWriteFailed), wErr))
	}

	fileName := file.FileType[strings.ToLower(entryType)]
	return s.toolOK(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPAddedFormat), entryType, fileName))
}

// toolComplete marks a task as done by number or text match.
func (s *Server) toolComplete(
	id json.RawMessage, args map[string]interface{},
) *Response {
	query, _ := args["query"].(string)
	if query == "" {
		return s.toolError(id, assets.TextDesc(assets.TextDescKeyMCPQueryRequired))
	}

	completedTask, err := complete.Task(query, s.contextDir)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolOK(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPCompletedFormat), completedTask))
}

// toolDrift runs drift detection and returns the report.
func (s *Server) toolDrift(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	report := drift.Detect(ctx)

	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPDriftStatusFormat), report.Status())

	if len(report.Violations) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftViolations))
		for _, v := range report.Violations {
			_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPDriftIssueFormat),
				v.Type, v.File, v.Message)
		}
		sb.WriteString(config.NewlineLF)
	}

	if len(report.Warnings) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftWarnings))
		for _, w := range report.Warnings {
			_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPDriftIssueFormat),
				w.Type, w.File, w.Message)
		}
		sb.WriteString(config.NewlineLF)
	}

	if len(report.Passed) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftPassed))
		for _, p := range report.Passed {
			_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPDriftPassedFormat), p)
		}
	}

	return s.toolOK(id, sb.String())
}

// toolOK builds a successful tool result.
func (s *Server) toolOK(id json.RawMessage, text string) *Response {
	return s.ok(id, CallToolResult{
		Content: []ToolContent{{Type: config.MCPContentTypeText, Text: text}},
	})
}

// toolError builds a tool error result.
func (s *Server) toolError(id json.RawMessage, msg string) *Response {
	return s.ok(id, CallToolResult{
		Content: []ToolContent{{Type: config.MCPContentTypeText, Text: msg}},
		IsError: true,
	})
}
