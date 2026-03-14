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
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/compact/core"
	remindcore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	taskcomplete "github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	entry2 "github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entry"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/task"
	"github.com/ActiveMemory/ctx/internal/validation"
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
	{
		Name:        mcp.MCPToolRecall,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolRecallDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				"limit": {
					Type:        mcp.MCPSchemaNumber,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropLimit),
				},
				"since": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropSince),
				},
			},
		},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name:        mcp.MCPToolWatchUpdate,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolWatchUpdateDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				"type": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropEntryType),
				},
				"content": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropMainContent),
				},
				"context": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropCtxBg),
				},
				"rationale": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropRationale),
				},
				"consequences": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropConseq),
				},
				"lesson": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropLesson),
				},
				"application": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropApplication),
				},
			},
			Required: []string{"type", "content"},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name:        mcp.MCPToolCompact,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolCompactDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				"archive": {
					Type:        mcp.MCPSchemaBoolean,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropArchive),
				},
			},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name:        mcp.MCPToolNext,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolNextDesc),
		InputSchema: InputSchema{Type: mcp.MCPSchemaObject},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name:        mcp.MCPToolCheckTaskCompletion,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolCheckTaskDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				"recent_action": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropRecentAct),
				},
			},
		},
		Annotations: &ToolAnnotations{ReadOnlyHint: true},
	},
	{
		Name:        mcp.MCPToolSessionEvent,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolSessionDesc),
		InputSchema: InputSchema{
			Type: mcp.MCPSchemaObject,
			Properties: map[string]Property{
				"type": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropEventType),
				},
				"caller": {
					Type:        mcp.MCPSchemaString,
					Description: assets.TextDesc(assets.TextDescKeyMCPToolPropCaller),
				},
			},
			Required: []string{"type"},
		},
		Annotations: &ToolAnnotations{},
	},
	{
		Name:        mcp.MCPToolRemind,
		Description: assets.TextDesc(assets.TextDescKeyMCPToolRemindDesc),
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

	s.session.recordToolCall()

	switch params.Name {
	case mcp.MCPToolStatus:
		return s.toolStatus(req.ID)
	case mcp.MCPToolAdd:
		return s.toolAdd(req.ID, params.Arguments)
	case mcp.MCPToolComplete:
		return s.toolComplete(req.ID, params.Arguments)
	case mcp.MCPToolDrift:
		return s.toolDrift(req.ID)
	case mcp.MCPToolRecall:
		return s.toolRecall(req.ID, params.Arguments)
	case mcp.MCPToolWatchUpdate:
		return s.toolWatchUpdate(req.ID, params.Arguments)
	case mcp.MCPToolCompact:
		return s.toolCompact(req.ID, params.Arguments)
	case mcp.MCPToolNext:
		return s.toolNext(req.ID)
	case mcp.MCPToolCheckTaskCompletion:
		return s.toolCheckTaskCompletion(req.ID, params.Arguments)
	case mcp.MCPToolSessionEvent:
		return s.toolSessionEvent(req.ID, params.Arguments)
	case mcp.MCPToolRemind:
		return s.toolRemind(req.ID)
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
	if err := validation.ValidateBoundary(s.contextDir); err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPBoundaryViolation), err))
	}

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
	if err := validation.ValidateBoundary(s.contextDir); err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPBoundaryViolation), err))
	}

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

// toolRecall queries recent session history.
func (s *Server) toolRecall(
	id json.RawMessage, args map[string]interface{},
) *Response {
	limit := 5
	if v, ok := args["limit"].(float64); ok && v > 0 {
		limit = int(v)
	}

	var sinceFilter time.Time
	if v, ok := args["since"].(string); ok && v != "" {
		parsed, parseErr := time.Parse("2006-01-02", v)
		if parseErr != nil {
			return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPInvalidSinceDate), parseErr))
		}
		sinceFilter = parsed
	}

	sessions, err := parser.FindSessions()
	if err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPFindSessionsFailed), err))
	}

	// Apply since filter.
	if !sinceFilter.IsZero() {
		var filtered []*parser.Session
		for _, sess := range sessions {
			if sess.StartTime.After(sinceFilter) || sess.StartTime.Equal(sinceFilter) {
				filtered = append(filtered, sess)
			}
		}
		sessions = filtered
	}

	// Apply limit.
	if len(sessions) > limit {
		sessions = sessions[:limit]
	}

	if len(sessions) == 0 {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPNoSessions))
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPSessionsFoundFormat)+"%s%s", len(sessions), token.NewlineLF, token.NewlineLF)

	for i, sess := range sessions {
		duration := sess.Duration.Round(time.Second)
		fmt.Fprintf(&sb, "%d. %s", i+1, sess.StartTime.Format("2006-01-02 15:04"))
		if sess.Project != "" {
			fmt.Fprintf(&sb, " [%s]", sess.Project)
		}
		fmt.Fprintf(&sb, " (%s, %d turns)", duration, sess.TurnCount)
		sb.WriteString(token.NewlineLF)

		if sess.FirstUserMsg != "" {
			fmt.Fprintf(&sb, "   %s", sess.FirstUserMsg)
			sb.WriteString(token.NewlineLF)
		}
	}

	return s.toolOK(id, sb.String())
}

// toolWatchUpdate applies a structured context-update to .context/ files.
func (s *Server) toolWatchUpdate(
	id json.RawMessage, args map[string]interface{},
) *Response {
	if err := validation.ValidateBoundary(s.contextDir); err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPBoundaryViolation), err))
	}

	entryType, _ := args["type"].(string)
	content, _ := args["content"].(string)

	if entryType == "" || content == "" {
		return s.toolError(id, assets.TextDesc(assets.TextDescKeyMCPTypeContentRequired))
	}

	// Handle "complete" type as a special case — delegate to ctx_complete.
	if entryType == "complete" {
		completedTask, err := taskcomplete.CompleteTask(content, s.contextDir)
		if err != nil {
			return s.toolError(id, err.Error())
		}
		s.session.queuePendingUpdate(PendingUpdate{
			Type:     entryType,
			Content:  content,
			QueuedAt: time.Now(),
		})
		return s.toolOK(id,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPWatchCompletedFormat), completedTask)+token.NewlineLF+
				assets.TextDesc(assets.TextDescKeyMCPReviewStatus))
	}

	params := entry.Params{
		Type:       entryType,
		Content:    content,
		ContextDir: s.contextDir,
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

	if vErr := entry.Validate(params, nil); vErr != nil {
		return s.toolError(id, vErr.Error())
	}

	if wErr := entry.Write(params); wErr != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPWriteFailed), wErr))
	}

	fileName := entry2.ToCtxFile[strings.ToLower(entryType)]
	s.session.recordAdd(entryType)
	s.session.queuePendingUpdate(PendingUpdate{
		Type:    entryType,
		Content: content,
		Attrs: map[string]string{
			"file": fileName,
		},
		QueuedAt: time.Now(),
	})

	return s.toolOK(id,
		fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPWroteFormat), entryType, fileName)+token.NewlineLF+
			assets.TextDesc(assets.TextDescKeyMCPReviewStatus))
}

// toolCompact moves completed tasks to the archive section.
func (s *Server) toolCompact(
	id json.RawMessage, args map[string]interface{},
) *Response {
	if err := validation.ValidateBoundary(s.contextDir); err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPBoundaryViolation), err))
	}

	archive := false
	if v, ok := args["archive"].(bool); ok {
		archive = v
	}

	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	var sb strings.Builder
	changes := 0

	// Process TASKS.md.
	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile != nil {
		content := string(tasksFile.Content)
		lines := strings.Split(content, token.NewlineLF)

		blocks := core.ParseTaskBlocks(lines)

		var archivableBlocks []core.TaskBlock
		for _, block := range blocks {
			if block.IsArchivable {
				archivableBlocks = append(archivableBlocks, block)
				fmt.Fprintf(&sb, "Moved: %s%s",
					core.TruncateString(block.ParentTaskText(), 50), token.NewlineLF)
			}
		}

		if len(archivableBlocks) > 0 {
			newLines := core.RemoveBlocksFromLines(lines, archivableBlocks)

			// Add blocks to the Completed section.
			for i, line := range newLines {
				if strings.HasPrefix(line, assets.HeadingCompleted) {
					insertIdx := i + 1
					for insertIdx < len(newLines) && newLines[insertIdx] != "" &&
						!strings.HasPrefix(newLines[insertIdx], token.HeadingLevelTwoStart) {
						insertIdx++
					}

					var blocksToInsert []string
					for _, block := range archivableBlocks {
						blocksToInsert = append(blocksToInsert, block.Lines...)
					}

					newLines = append(newLines[:insertIdx],
						append(blocksToInsert, newLines[insertIdx:]...)...)
					break
				}
			}

			newContent := strings.Join(newLines, token.NewlineLF)
			if newContent != content {
				if writeErr := writeContextFile(tasksFile.Path, []byte(newContent)); writeErr != nil {
					return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPWriteFailed), writeErr))
				}
			}
			changes += len(archivableBlocks)
		}

		// Archive old tasks if requested.
		if archive && len(archivableBlocks) > 0 {
			var archiveContent string
			for _, block := range archivableBlocks {
				archiveContent += block.BlockContent() + token.NewlineLF + token.NewlineLF
			}
			if _, archiveErr := core.WriteArchive("tasks", assets.HeadingArchivedTasks, archiveContent); archiveErr != nil {
				fmt.Fprintf(&sb, "Archive warning: %v%s", archiveErr, token.NewlineLF)
			}
		}
	}

	// Process other files for empty sections.
	for _, f := range ctx.Files {
		if f.Name == ctxCfg.Task {
			continue
		}
		cleaned, count := core.RemoveEmptySections(string(f.Content))
		if count > 0 {
			if writeErr := writeContextFile(f.Path, []byte(cleaned)); writeErr == nil {
				fmt.Fprintf(&sb, "Removed %d empty sections from %s%s",
					count, f.Name, token.NewlineLF)
				changes += count
			}
		}
	}

	if changes == 0 {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPCompactClean))
	}

	fmt.Fprintf(&sb, "%s"+assets.TextDesc(assets.TextDescKeyMCPCompactedFormat)+"%s",
		token.NewlineLF, changes, token.NewlineLF)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPReviewStatus))

	return s.toolOK(id, sb.String())
}

// toolNext suggests the next pending task.
func (s *Server) toolNext(id json.RawMessage) *Response {
	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPNoTasks))
	}

	content := string(tasksFile.Content)
	lines := strings.Split(content, token.NewlineLF)

	// Find the first pending top-level task.
	inCompletedSection := false
	pendingIdx := 0

	for _, line := range lines {
		if strings.HasPrefix(line, assets.HeadingCompleted) {
			inCompletedSection = true
			continue
		}
		if strings.HasPrefix(line, token.HeadingLevelTwoStart) && inCompletedSection {
			inCompletedSection = false
		}
		if inCompletedSection {
			continue
		}

		match := regex.Task.FindStringSubmatch(line)
		if match == nil || !task.Pending(match) {
			continue
		}

		// Skip subtasks.
		if task.SubTask(match) {
			continue
		}

		pendingIdx++
		return s.toolOK(id, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPNextTaskFormat), pendingIdx, task.Content(match)))
	}

	return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPAllTasksComplete))
}

// toolCheckTaskCompletion checks if a recent action completed any pending tasks.
func (s *Server) toolCheckTaskCompletion(
	id json.RawMessage, args map[string]interface{},
) *Response {
	recentAction, _ := args["recent_action"].(string)

	ctx, err := context.Load(s.contextDir)
	if err != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPLoadContext), err))
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return s.toolOK(id, "")
	}

	content := string(tasksFile.Content)
	lines := strings.Split(content, token.NewlineLF)

	inCompletedSection := false
	taskNum := 0

	for _, line := range lines {
		if strings.HasPrefix(line, assets.HeadingCompleted) {
			inCompletedSection = true
			continue
		}
		if strings.HasPrefix(line, token.HeadingLevelTwoStart) && inCompletedSection {
			inCompletedSection = false
		}
		if inCompletedSection {
			continue
		}

		match := regex.Task.FindStringSubmatch(line)
		if match == nil || !task.Pending(match) {
			continue
		}
		if task.SubTask(match) {
			continue
		}

		taskNum++
		taskText := task.Content(match)

		// Check for keyword overlap between the recent action and the task.
		if recentAction != "" && containsOverlap(recentAction, taskText) {
			return s.toolOK(id, fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPCheckTaskFormat)+token.NewlineLF+
					assets.TextDesc(assets.TextDescKeyMCPCheckTaskHint), taskNum, taskText, taskNum))
		}
	}

	return s.toolOK(id, "")
}

// toolSessionEvent handles session lifecycle events.
func (s *Server) toolSessionEvent(
	id json.RawMessage, args map[string]interface{},
) *Response {
	eventType, _ := args["type"].(string)
	if eventType == "" {
		return s.toolError(id, assets.TextDesc(assets.TextDescKeyMCPEventTypeRequired))
	}

	switch eventType {
	case "start":
		s.session = newSessionState(s.contextDir)
		if caller, ok := args["caller"].(string); ok && caller != "" {
			return s.toolOK(id, fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPSessionStartedCallerFormat), caller, s.contextDir))
		}
		return s.toolOK(id, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPSessionStartedFormat), s.contextDir))

	case "end":
		pending := s.session.pendingCount()
		var sb strings.Builder
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPSessionEnding))
		sb.WriteString(token.NewlineLF)

		if pending > 0 {
			fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPPendingUpdatesFormat)+"%s",
				pending, token.NewlineLF)
			for i, pu := range s.session.pendingFlush {
				fmt.Fprintf(&sb, "  %d. [%s] %s%s",
					i+1, pu.Type, core.TruncateString(pu.Content, 60), token.NewlineLF)
			}
			sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPReviewPending))
		} else {
			sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPNoPending))
		}

		fmt.Fprintf(&sb, "%s"+assets.TextDesc(assets.TextDescKeyMCPSessionStatsFormat),
			token.NewlineLF, s.session.toolCalls, totalAdds(s.session.addsPerformed))

		return s.toolOK(id, sb.String())

	default:
		return s.toolError(id,
			fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPUnknownEventType), eventType))
	}
}

// toolRemind lists pending session-scoped reminders.
func (s *Server) toolRemind(id json.RawMessage) *Response {
	reminders, readErr := remindcore.ReadReminders()
	if readErr != nil {
		return s.toolError(id, fmt.Sprintf(assets.TextDesc(assets.TextDescKeyMCPReadRemindersFailed), readErr))
	}

	if len(reminders) == 0 {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPNoReminders))
	}

	today := time.Now().Format("2006-01-02")
	var sb strings.Builder
	fmt.Fprintf(&sb, assets.TextDesc(assets.TextDescKeyMCPRemindersFormat)+"%s", len(reminders), token.NewlineLF)

	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf(" (after %s, not yet due)", *r.After)
			}
		}
		fmt.Fprintf(&sb, "  [%d] %s%s%s", r.ID, r.Message, annotation, token.NewlineLF)
	}

	return s.toolOK(id, sb.String())
}

// containsOverlap checks if two strings share meaningful words.
func containsOverlap(action, taskText string) bool {
	actionLower := strings.ToLower(action)
	taskLower := strings.ToLower(taskText)

	// Split task text into words, check if any appear in the action.
	words := strings.Fields(taskLower)
	matchCount := 0
	for _, w := range words {
		if len(w) < 4 {
			continue // Skip short common words.
		}
		if strings.Contains(actionLower, w) {
			matchCount++
		}
	}

	// Require at least 2 word matches for a reasonable signal.
	return matchCount >= 2
}

// totalAdds sums all entry add counts.
func totalAdds(m map[string]int) int {
	total := 0
	for _, v := range m {
		total += v
	}
	return total
}

// writeContextFile writes content to a context file with standard permissions.
func writeContextFile(path string, data []byte) error {
	return os.WriteFile(path, data, fs.PermFile)
}
