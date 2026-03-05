//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/cli/add"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/task"
)

// toolDefs defines all available MCP tools.
var toolDefs = []Tool{
	{
		Name:        "ctx_status",
		Description: "Show context health: file count, token estimate, and file summaries",
		InputSchema: InputSchema{Type: "object"},
	},
	{
		Name:        "ctx_add",
		Description: "Add a task, decision, learning, or convention to the context",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"type": {
					Type:        "string",
					Description: "Entry type to add",
					Enum:        []string{"task", "decision", "learning", "convention"},
				},
				"content": {
					Type:        "string",
					Description: "Title or main content of the entry",
				},
				"priority": {
					Type:        "string",
					Description: "Priority level (for tasks only)",
					Enum:        []string{"high", "medium", "low"},
				},
				"context": {
					Type:        "string",
					Description: "Context field (required for decisions and learnings)",
				},
				"rationale": {
					Type:        "string",
					Description: "Rationale (required for decisions)",
				},
				"consequences": {
					Type:        "string",
					Description: "Consequences (required for decisions)",
				},
				"lesson": {
					Type:        "string",
					Description: "Lesson learned (required for learnings)",
				},
				"application": {
					Type:        "string",
					Description: "How to apply this lesson (required for learnings)",
				},
			},
			Required: []string{"type", "content"},
		},
	},
	{
		Name:        "ctx_complete",
		Description: "Mark a task as done by number or text match",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]Property{
				"query": {
					Type:        "string",
					Description: "Task number (e.g. '1') or search text to match",
				},
			},
			Required: []string{"query"},
		},
	},
	{
		Name:        "ctx_drift",
		Description: "Detect stale or invalid context: dead paths, missing files, staleness",
		InputSchema: InputSchema{Type: "object"},
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
		return s.error(req.ID, errCodeInvalidArg, "invalid params")
	}

	switch params.Name {
	case "ctx_status":
		return s.toolStatus(req.ID)
	case "ctx_add":
		return s.toolAdd(req.ID, params.Arguments)
	case "ctx_complete":
		return s.toolComplete(req.ID, params.Arguments)
	case "ctx_drift":
		return s.toolDrift(req.ID)
	default:
		return s.error(req.ID, errCodeNotFound,
			fmt.Sprintf("unknown tool: %s", params.Name))
	}
}

// toolStatus loads context and returns a status summary.
func (s *Server) toolStatus(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf("failed to load context: %v", err))
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Context: %s\n", ctx.Dir)
	fmt.Fprintf(&sb, "Files: %d\n", len(ctx.Files))
	fmt.Fprintf(&sb, "Tokens: ~%d\n\n", ctx.TotalTokens)

	for _, f := range ctx.Files {
		status := "OK"
		if f.IsEmpty {
			status = "EMPTY"
		}
		fmt.Fprintf(&sb, "  %-22s %6d tokens  [%s]\n",
			f.Name, f.Tokens, status)
	}

	return s.toolOK(id, sb.String())
}

// toolAdd adds an entry to a context file.
func (s *Server) toolAdd(
	id json.RawMessage, args map[string]interface{},
) *Response {
	entryType, _ := args["type"].(string)
	content, _ := args["content"].(string)

	if entryType == "" || content == "" {
		return s.toolError(id, "type and content are required")
	}

	params := add.EntryParams{
		Type:    entryType,
		Content: content,
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
	if vErr := add.ValidateEntry(params); vErr != nil {
		return s.toolError(id, vErr.Error())
	}

	if wErr := add.WriteEntry(params); wErr != nil {
		return s.toolError(id, fmt.Sprintf("write failed: %v", wErr))
	}

	fileName := config.FileType[strings.ToLower(entryType)]
	return s.toolOK(id, fmt.Sprintf("Added %s to %s", entryType, fileName))
}

// toolComplete marks a task as done by number or text match.
func (s *Server) toolComplete(
	id json.RawMessage, args map[string]interface{},
) *Response {
	query, _ := args["query"].(string)
	if query == "" {
		return s.toolError(id, "query is required")
	}

	filePath := filepath.Join(s.contextDir, config.FileTask)

	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return s.toolError(id, "TASKS.md not found")
	}

	lines := strings.Split(string(content), config.NewlineLF)

	var taskNumber int
	isNumber := false
	if num, parseErr := strconv.Atoi(query); parseErr == nil {
		taskNumber = num
		isNumber = true
	}

	currentTaskNum := 0
	matchedLine := -1
	matchedTask := ""

	for i, line := range lines {
		match := config.RegExTask.FindStringSubmatch(line)
		if match != nil && task.Pending(match) {
			currentTaskNum++
			taskText := task.Content(match)

			if isNumber && currentTaskNum == taskNumber {
				matchedLine = i
				matchedTask = taskText
				break
			}

			if !isNumber && strings.Contains(
				strings.ToLower(taskText), strings.ToLower(query),
			) {
				if matchedLine != -1 {
					return s.toolError(id,
						fmt.Sprintf("multiple tasks match %q; be more specific or use task number", query))
				}
				matchedLine = i
				matchedTask = taskText
			}
		}
	}

	if matchedLine == -1 {
		return s.toolError(id, fmt.Sprintf("no task matching %q found", query))
	}

	lines[matchedLine] = config.RegExTask.ReplaceAllString(
		lines[matchedLine], "$1- [x] $3",
	)

	newContent := strings.Join(lines, config.NewlineLF)
	if wErr := os.WriteFile(filePath, []byte(newContent), config.PermFile); wErr != nil {
		return s.toolError(id, fmt.Sprintf("failed to write TASKS.md: %v", wErr))
	}

	return s.toolOK(id, fmt.Sprintf("Completed: %s", matchedTask))
}

// toolDrift runs drift detection and returns the report.
func (s *Server) toolDrift(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf("failed to load context: %v", err))
	}

	report := drift.Detect(ctx)

	var sb strings.Builder
	fmt.Fprintf(&sb, "Status: %s\n\n", report.Status())

	if len(report.Violations) > 0 {
		sb.WriteString("Violations:\n")
		for _, v := range report.Violations {
			fmt.Fprintf(&sb, "  - [%s] %s: %s\n",
				v.Type, v.File, v.Message)
		}
		sb.WriteString("\n")
	}

	if len(report.Warnings) > 0 {
		sb.WriteString("Warnings:\n")
		for _, w := range report.Warnings {
			fmt.Fprintf(&sb, "  - [%s] %s: %s\n",
				w.Type, w.File, w.Message)
		}
		sb.WriteString("\n")
	}

	if len(report.Passed) > 0 {
		sb.WriteString("Passed:\n")
		for _, p := range report.Passed {
			fmt.Fprintf(&sb, "  - %s\n", p)
		}
	}

	return s.toolOK(id, sb.String())
}

// toolOK builds a successful tool result.
func (s *Server) toolOK(id json.RawMessage, text string) *Response {
	return s.ok(id, CallToolResult{
		Content: []ToolContent{{Type: "text", Text: text}},
	})
}

// toolError builds a tool error result.
func (s *Server) toolError(id json.RawMessage, msg string) *Response {
	return s.ok(id, CallToolResult{
		Content: []ToolContent{{Type: "text", Text: msg}},
		IsError: true,
	})
}
