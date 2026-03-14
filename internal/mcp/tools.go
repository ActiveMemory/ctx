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
	taskcomplete "github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	entry2 "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entry"
)

// toolDefs defines all available MCP tools.
var toolDefs = []Tool{
	{
		Name:        mcp.MCPToolStatus,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolStatusDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name:        mcp.MCPToolAdd,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolAddDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				cli.AttrType: {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropType),
					Enum:        []string{"task", "decision", "learning", "convention"},
				},
				"content": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropContent),
				},
				"priority": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropPriority),
					Enum:        []string{"high", "medium", "low"},
				},
				cli.AttrContext: {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropContext),
				},
				cli.AttrRationale: {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropRationale),
				},
				cli.AttrConsequences: {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropConseq),
				},
				cli.AttrLesson: {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropLesson),
				},
				cli.AttrApplication: {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropApplication),
				},
			},
			Required: []string{cli.AttrType, "content"},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name:        mcp.MCPToolComplete,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolCompleteDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				"query": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropQuery),
				},
			},
			Required: []string{"query"},
		},
		Annotations: &ToolAnnotations{IdempotentHint: true},
	},
	{
		Name:        mcp.MCPToolDrift,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolDriftDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
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
	case mcp.MCPToolStatus:
		return s.toolStatus(req.ID)
	case mcp.MCPToolAdd:
		return s.toolAdd(req.ID, params.Arguments)
	case mcp.MCPToolComplete:
		return s.toolComplete(req.ID, params.Arguments)
	case mcp.MCPToolDrift:
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
	entryType, _ := args[cli.AttrType].(string)
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

	fileName := entry2.ToCtxFile[strings.ToLower(entryType)]
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

	completedTask, err := taskcomplete.CompleteTask(query, s.contextDir)
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
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Warnings) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftWarnings))
		for _, w := range report.Warnings {
			_, _ = fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPDriftIssueFormat),
				w.Type, w.File, w.Message)
		}
		sb.WriteString(token.NewlineLF)
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
		Content: []ToolContent{{Type: mcp.MCPContentTypeText, Text: text}},
	})
}

// toolError builds a tool error result.
func (s *Server) toolError(id json.RawMessage, msg string) *Response {
	return s.ok(id, CallToolResult{
		Content: []ToolContent{{Type: mcp.MCPContentTypeText, Text: msg}},
		IsError: true,
	})
}
