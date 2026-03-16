//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets"
	remindcore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	taskcomplete "github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	archiveCfg "github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	entryCfg "github.com/ActiveMemory/ctx/internal/config/entry"
	configfs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/event"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	timeCfg "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entry"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/session"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/tidy"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// checkBoundary validates the context directory boundary and returns
// an error response if the check fails; nil otherwise.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *proto.Response: non-nil error response if boundary is invalid
func (s *Server) checkBoundary(id json.RawMessage) *proto.Response {
	if bErr := validation.ValidateBoundary(s.contextDir); bErr != nil {
		return s.toolError(
			id,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPBoundaryViolation), bErr),
		)
	}
	return nil
}

// loadContext loads the context directory and returns an error
// response if loading fails; nil otherwise.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *context.Context: loaded context (nil on error)
//   - *proto.Response: non-nil error response if load fails
func (s *Server) loadContext(
	id json.RawMessage,
) (*context.Context, *proto.Response) {
	ctx, loadErr := context.Load(s.contextDir)
	if loadErr != nil {
		return nil, s.toolError(
			id,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPLoadContext), loadErr),
		)
	}
	return ctx, nil
}

// extractEntryArgs validates the context directory boundary and
// extracts the required type/content pair from MCP tool arguments.
//
// Parameters:
//   - id: JSON-RPC request ID (for error responses)
//   - args: tool arguments
//
// Returns:
//   - entryType: extracted entry type string
//   - content: extracted content string
//   - errResp: non-nil error response if boundary or presence
//     checks fail; caller should return it immediately
func (s *Server) extractEntryArgs(
	id json.RawMessage, args map[string]interface{},
) (entryType, content string, errResp *proto.Response) {
	if errResp = s.checkBoundary(id); errResp != nil {
		return "", "", errResp
	}

	entryType, _ = args[cli.AttrType].(string)
	content, _ = args[field.Content].(string)

	if entryType == "" || content == "" {
		return "", "", s.toolError(
			id, assets.TextDesc(assets.TextDescKeyMCPTypeContentRequired),
		)
	}

	return entryType, content, nil
}

// validateAndWriteEntry validates entry params, writes the entry,
// and returns the target context filename.
//
// Parameters:
//   - id: JSON-RPC request ID (for error responses)
//   - params: populated entry parameters
//
// Returns:
//   - fileName: the context file that was written to
//   - errResp: non-nil error response if validation or write fails
func (s *Server) validateAndWriteEntry(
	id json.RawMessage, params entry.Params,
) (fileName string, errResp *proto.Response) {
	if vErr := entry.Validate(params, nil); vErr != nil {
		return "", s.toolError(id, vErr.Error())
	}

	if wErr := entry.Write(params); wErr != nil {
		return "", s.toolError(
			id, fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPWriteFailed), wErr),
		)
	}

	return entryCfg.ToCtxFile[strings.ToLower(params.Type)], nil
}

// handleToolsList returns all available MCP tools.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *Response: tool list result
func (s *Server) handleToolsList(req proto.Request) *proto.Response {
	return s.ok(req.ID, proto.ToolListResult{Tools: proto.ToolDefs})
}

// handleToolsCall dispatches a tool call to the appropriate handler.
//
// Parameters:
//   - req: the MCP request containing tool name and arguments
//
// Returns:
//   - *Response: tool result or error
func (s *Server) handleToolsCall(req proto.Request) *proto.Response {
	var params proto.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(
			req.ID, proto.ErrCodeInvalidArg,
			assets.TextDesc(assets.TextDescKeyMCPInvalidParams),
		)
	}

	s.session.RecordToolCall()

	switch params.Name {
	case tool.Status:
		return s.toolStatus(req.ID)
	case tool.Add:
		return s.toolAdd(req.ID, params.Arguments)
	case tool.Complete:
		return s.toolComplete(req.ID, params.Arguments)
	case tool.Drift:
		return s.toolDrift(req.ID)
	case tool.Recall:
		return s.toolRecall(req.ID, params.Arguments)
	case tool.WatchUpdate:
		return s.toolWatchUpdate(req.ID, params.Arguments)
	case tool.Compact:
		return s.toolCompact(req.ID, params.Arguments)
	case tool.Next:
		return s.toolNext(req.ID)
	case tool.CheckTaskCompletion:
		return s.toolCheckTaskCompletion(req.ID, params.Arguments)
	case tool.SessionEvent:
		return s.toolSessionEvent(req.ID, params.Arguments)
	case tool.Remind:
		return s.toolRemind(req.ID)
	default:
		return s.error(
			req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPUnknownTool),
				params.Name,
			),
		)
	}
}

// toolStatus loads context and returns a status summary.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: context status with file list and token counts
func (s *Server) toolStatus(id json.RawMessage) *proto.Response {
	ctx, errResp := s.loadContext(id)
	if errResp != nil {
		return errResp
	}

	var sb strings.Builder
	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPStatusContextFormat), ctx.Dir,
	)
	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPStatusFilesFormat), len(ctx.Files),
	)
	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPStatusTokensFormat), ctx.TotalTokens,
	)

	for _, f := range ctx.Files {
		status := assets.TextDesc(assets.TextDescKeyMCPStatusOK)
		if f.IsEmpty {
			status = assets.TextDesc(assets.TextDescKeyMCPStatusEmpty)
		}
		_, _ = fmt.Fprintf(
			&sb, assets.TextDesc(assets.TextDescKeyMCPStatusFileFormat),
			f.Name, f.Tokens, status,
		)
	}

	return s.toolOK(id, sb.String())
}

// toolAdd adds an entry to a context file.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (type, content, optional fields)
//
// Returns:
//   - *Response: confirmation or validation error
func (s *Server) toolAdd(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	entryType, content, errResp := s.extractEntryArgs(id, args)
	if errResp != nil {
		return errResp
	}

	params := entry.Params{
		Type:       entryType,
		Content:    content,
		ContextDir: s.contextDir,
	}

	applyOptionalFields(&params, args)

	fileName, errResp := s.validateAndWriteEntry(id, params)
	if errResp != nil {
		return errResp
	}

	return s.toolOK(
		id,
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPAddedFormat),
			entryType, fileName,
		),
	)
}

// toolComplete marks a task as done by number or text match.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (query)
//
// Returns:
//   - *Response: completed task name or error
func (s *Server) toolComplete(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	if errResp := s.checkBoundary(id); errResp != nil {
		return errResp
	}

	query, _ := args[field.Query].(string)
	if query == "" {
		return s.toolError(id, assets.TextDesc(assets.TextDescKeyMCPQueryRequired))
	}

	completedTask, err := taskcomplete.CompleteTask(query, s.contextDir)
	if err != nil {
		return s.toolError(id, err.Error())
	}

	return s.toolOK(
		id, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPCompletedFormat), completedTask),
	)
}

// toolDrift runs drift detection and returns the report.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: drift report with violations and warnings
func (s *Server) toolDrift(id json.RawMessage) *proto.Response {
	ctx, errResp := s.loadContext(id)
	if errResp != nil {
		return errResp
	}

	report := drift.Detect(ctx)

	var sb strings.Builder
	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPDriftStatusFormat),
		report.Status(),
	)

	if len(report.Violations) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftViolations))
		for _, v := range report.Violations {
			_, _ = fmt.Fprintf(
				&sb, assets.TextDesc(assets.TextDescKeyMCPDriftIssueFormat),
				v.Type, v.File, v.Message,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Warnings) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftWarnings))
		for _, w := range report.Warnings {
			_, _ = fmt.Fprintf(
				&sb, assets.TextDesc(assets.TextDescKeyMCPDriftIssueFormat),
				w.Type, w.File, w.Message,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Passed) > 0 {
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPDriftPassed))
		for _, p := range report.Passed {
			_, _ = fmt.Fprintf(
				&sb, assets.TextDesc(assets.TextDescKeyMCPDriftPassedFormat), p,
			)
		}
	}

	return s.toolOK(id, sb.String())
}

// toolOK builds a successful tool result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - text: result text content
//
// Returns:
//   - *Response: success response with text content
func (s *Server) toolOK(id json.RawMessage, text string) *proto.Response {
	return s.ok(
		id,
		proto.CallToolResult{
			Content: []proto.ToolContent{
				{Type: mime.ContentTypeText, Text: text},
			},
		})
}

// toolError builds a tool error result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - msg: error message
//
// Returns:
//   - *Response: error response with IsError flag
func (s *Server) toolError(id json.RawMessage, msg string) *proto.Response {
	return s.ok(id, proto.CallToolResult{
		Content: []proto.ToolContent{{Type: mime.ContentTypeText, Text: msg}},
		IsError: true,
	})
}

// toolRecall queries recent session history.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (limit, since)
//
// Returns:
//   - *Response: formatted session list or empty result
func (s *Server) toolRecall(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	limit := cfg.DefaultRecallLimit
	if v, ok := args[field.Limit].(float64); ok && v > 0 {
		limit = int(v)
	}

	var sinceFilter time.Time
	if v, ok := args[field.Since].(string); ok && v != "" {
		parsed, parseErr := time.Parse(timeCfg.DateFormat, v)
		if parseErr != nil {
			return s.toolError(
				id, fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPInvalidSinceDate), parseErr),
			)
		}
		sinceFilter = parsed
	}

	sessions, err := parser.FindSessions()
	if err != nil {
		return s.toolError(
			id, fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPFindSessionsFailed), err),
		)
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
	_, _ = fmt.Fprintf(&sb,
		assets.TextDesc(assets.TextDescKeyMCPSessionsFoundFormat),
		len(sessions),
	)

	for i, sess := range sessions {
		duration := sess.Duration.Round(time.Second)
		_, _ = fmt.Fprintf(
			&sb,
			assets.TextDesc(assets.TextDescKeyMCPRecallItemFormat),
			i+1, sess.StartTime.Format(timeCfg.DateTimeFormat),
		)
		if sess.Project != "" {
			_, _ = fmt.Fprintf(
				&sb, assets.TextDesc(assets.TextDescKeyMCPRecallProjectFormat),
				sess.Project,
			)
		}
		_, _ = fmt.Fprintf(
			&sb, assets.TextDesc(assets.TextDescKeyMCPRecallDurationFormat),
			duration, sess.TurnCount,
		)
		sb.WriteString(token.NewlineLF)

		if sess.FirstUserMsg != "" {
			_, _ = fmt.Fprintf(
				&sb, assets.TextDesc(assets.TextDescKeyMCPRecallFirstMsgFormat),
				sess.FirstUserMsg,
			)
			sb.WriteString(token.NewlineLF)
		}
	}

	return s.toolOK(id, sb.String())
}

// toolWatchUpdate applies a structured context-update to .context/ files.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (type, content, optional fields)
//
// Returns:
//   - *Response: confirmation with file name or error
func (s *Server) toolWatchUpdate(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	entryType, content, errResp := s.extractEntryArgs(id, args)
	if errResp != nil {
		return errResp
	}

	// Handle the "complete" type as a special case: delegate to ctx_complete.
	if entryType == entryCfg.Complete {
		completedTask, completeErr := taskcomplete.CompleteTask(
			content, s.contextDir)
		if completeErr != nil {
			return s.toolError(id, completeErr.Error())
		}
		s.session.QueuePendingUpdate(session.PendingUpdate{
			Type:     entryType,
			Content:  content,
			QueuedAt: time.Now(),
		})
		return s.toolOK(
			id,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPWatchCompletedFormat),
				completedTask,
			)+token.NewlineLF+
				assets.TextDesc(assets.TextDescKeyMCPReviewStatus),
		)
	}

	params := entry.Params{
		Type:       entryType,
		Content:    content,
		ContextDir: s.contextDir,
	}

	applyOptionalFields(&params, args)

	fileName, errResp := s.validateAndWriteEntry(id, params)
	if errResp != nil {
		return errResp
	}

	s.session.RecordAdd(entryType)
	s.session.QueuePendingUpdate(session.PendingUpdate{
		Type:    entryType,
		Content: content,
		Attrs: map[string]string{
			field.AttrFile: fileName,
		},
		QueuedAt: time.Now(),
	})

	return s.toolOK(
		id,
		fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPWroteFormat),
			entryType, fileName,
		)+token.NewlineLF+
			assets.TextDesc(assets.TextDescKeyMCPReviewStatus),
	)
}

// toolCompact moves completed tasks to the archive section.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (archive)
//
// Returns:
//   - *Response: summary of compacted items or clean status
func (s *Server) toolCompact(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	if errResp := s.checkBoundary(id); errResp != nil {
		return errResp
	}

	archive := false
	if v, ok := args[field.Archive].(bool); ok {
		archive = v
	}

	ctx, errResp := s.loadContext(id)
	if errResp != nil {
		return errResp
	}

	result := tidy.CompactContext(ctx)

	// Write TASKS.md changes.
	if result.TasksFileUpdate != nil {
		if writeErr := os.WriteFile(
			result.TasksFileUpdate.Path,
			result.TasksFileUpdate.Content,
			configfs.PermFile,
		); writeErr != nil {
			return s.toolError(
				id,
				fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPWriteFailed), writeErr,
				),
			)
		}
	}

	// Write section-cleaned files.
	for _, fu := range result.SectionFileUpdates {
		if writeErr := os.WriteFile(
			fu.Path, fu.Content, configfs.PermFile,
		); writeErr != nil {
			return s.toolError(
				id,
				fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPWriteFailed), writeErr,
				),
			)
		}
	}

	// Archive old tasks if requested.
	var sb strings.Builder
	if archive && len(result.ArchivableBlocks) > 0 {
		var archiveContent string
		for _, block := range result.ArchivableBlocks {
			archiveContent += block.BlockContent() +
				token.NewlineLF + token.NewlineLF
		}
		if _, archiveErr := tidy.WriteArchive(
			archiveCfg.ArchiveScopeTasks,
			assets.HeadingArchivedTasks,
			archiveContent,
		); archiveErr != nil {
			_, _ = fmt.Fprintf(
				&sb,
				assets.TextDesc(assets.TextDescKeyMCPCompactArchiveWarning)+
					token.NewlineLF,
				archiveErr,
			)
		}
	}

	// Build response text.
	for _, taskText := range result.TasksMoved {
		_, _ = fmt.Fprintf(&sb,
			assets.TextDesc(
				assets.TextDescKeyMCPCompactMovedFormat)+token.NewlineLF,
			tidy.TruncateString(taskText, cfg.TruncateLen),
		)
	}
	for _, sc := range result.SectionsCleaned {
		_, _ = fmt.Fprintf(
			&sb,
			assets.TextDesc(assets.TextDescKeyMCPCompactRemovedSectFmt)+
				token.NewlineLF,
			sc.Removed, sc.FileName,
		)
	}

	if result.TotalChanges() == 0 {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPCompactClean))
	}

	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPCompactedFormat),
		result.TotalChanges(),
	)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPReviewStatus))

	return s.toolOK(id, sb.String())
}

// toolNext suggests the next pending task.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: next task or all-complete message
func (s *Server) toolNext(id json.RawMessage) *proto.Response {
	ctx, errResp := s.loadContext(id)
	if errResp != nil {
		return errResp
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPNoTasks))
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result *proto.Response
	eachPendingTask(lines, func(pt pendingTask) bool {
		result = s.toolOK(
			id,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPNextTaskFormat),
				pt.Index, pt.Content,
			),
		)
		return true // stop after first
	})

	if result != nil {
		return result
	}

	return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPAllTasksComplete))
}

// toolCheckTaskCompletion checks if a recent action completed
// any pending tasks.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (recent_action)
//
// Returns:
//   - *Response: nudge text if match found, empty otherwise
func (s *Server) toolCheckTaskCompletion(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	recentAction, _ := args[field.RecentAction].(string)

	ctx, errResp := s.loadContext(id)
	if errResp != nil {
		return errResp
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return s.toolOK(id, "")
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result *proto.Response
	eachPendingTask(lines, func(pt pendingTask) bool {
		if recentAction != "" && containsOverlap(recentAction, pt.Content) {
			result = s.toolOK(
				id,
				fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPCheckTaskFormat)+
						token.NewlineLF+
						assets.TextDesc(assets.TextDescKeyMCPCheckTaskHint),
					pt.Index, pt.Content, pt.Index,
				),
			)
			return true
		}
		return false
	})

	if result != nil {
		return result
	}

	return s.toolOK(id, "")
}

// toolSessionEvent handles session lifecycle events.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: tool arguments (type, caller)
//
// Returns:
//   - *Response: session status with stats on end
func (s *Server) toolSessionEvent(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	eventType, _ := args[cli.AttrType].(string)
	if eventType == "" {
		return s.toolError(id, assets.TextDesc(
			assets.TextDescKeyMCPEventTypeRequired),
		)
	}

	switch eventType {
	case event.Start:
		s.session = session.NewState(s.contextDir)
		if caller, ok := args[field.Caller].(string); ok && caller != "" {
			return s.toolOK(
				id,
				fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPSessionStartedCallerFormat),
					caller, s.contextDir,
				),
			)
		}
		return s.toolOK(
			id,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPSessionStartedFormat),
				s.contextDir,
			),
		)

	case event.End:
		pending := s.session.PendingCount()
		var sb strings.Builder
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPSessionEnding))
		sb.WriteString(token.NewlineLF)

		if pending > 0 {
			_, _ = fmt.Fprintf(
				&sb,
				assets.TextDesc(assets.TextDescKeyMCPPendingUpdatesFormat),
				pending,
			)
			for i, pu := range s.session.PendingFlush {
				_, _ = fmt.Fprintf(
					&sb,
					assets.TextDesc(assets.TextDescKeyMCPPendingItemFormat)+
						token.NewlineLF,
					i+1, pu.Type,
					tidy.TruncateString(pu.Content, cfg.TruncateContentLen),
				)
			}
			sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPReviewPending))
		} else {
			sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPNoPending))
		}

		_, _ = fmt.Fprintf(&sb,
			assets.TextDesc(assets.TextDescKeyMCPSessionStatsFormat),
			s.session.ToolCalls, totalAdds(s.session.AddsPerformed))

		return s.toolOK(id, sb.String())

	default:
		return s.toolError(id,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPUnknownEventType), eventType),
		)
	}
}

// toolRemind lists pending session-scoped reminders.
//
// Parameters:
//   - id: JSON-RPC request ID
//
// Returns:
//   - *Response: formatted reminder list or empty message
func (s *Server) toolRemind(id json.RawMessage) *proto.Response {
	reminders, readErr := remindcore.ReadReminders()
	if readErr != nil {
		return s.toolError(id,
			fmt.Sprintf(assets.TextDesc(
				assets.TextDescKeyMCPReadRemindersFailed), readErr),
		)
	}

	if len(reminders) == 0 {
		return s.toolOK(id, assets.TextDesc(assets.TextDescKeyMCPNoReminders))
	}

	today := time.Now().Format(timeCfg.DateFormat)
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb,
		assets.TextDesc(assets.TextDescKeyMCPRemindersFormat),
		len(reminders),
	)

	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPReminderNotDueFormat), *r.After,
				)
			}
		}
		_, _ = fmt.Fprintf(&sb, assets.TextDesc(
			assets.TextDescKeyMCPReminderItemFormat)+token.NewlineLF,
			r.ID, r.Message, annotation)
	}

	return s.toolOK(id, sb.String())
}
