//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	remindcore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	taskcomplete "github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	archiveCfg "github.com/ActiveMemory/ctx/internal/config/archive"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	entryCfg "github.com/ActiveMemory/ctx/internal/config/entry"
	configfs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/event"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	timeCfg "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entry"
	mcpErr "github.com/ActiveMemory/ctx/internal/err/mcp"
	"github.com/ActiveMemory/ctx/internal/mcp/handler/task"
	"github.com/ActiveMemory/ctx/internal/mcp/server/stat"
	"github.com/ActiveMemory/ctx/internal/mcp/session"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/tidy"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// Status loads context and returns a status summary.
//
// Returns:
//   - string: formatted status text with file list and token counts
//   - error: context load error
func (h *Handler) Status() (string, error) {
	ctx, loadErr := load.Do(h.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	var sb strings.Builder
	_, _ = fmt.Fprintf(
		&sb,
		desc.TextDesc(text.TextDescKeyMCPStatusContextFormat), ctx.Dir,
	)
	_, _ = fmt.Fprintf(
		&sb,
		desc.TextDesc(text.TextDescKeyMCPStatusFilesFormat), len(ctx.Files),
	)
	_, _ = fmt.Fprintf(
		&sb,
		desc.TextDesc(text.TextDescKeyMCPStatusTokensFormat), ctx.TotalTokens,
	)

	for _, f := range ctx.Files {
		status := desc.TextDesc(text.TextDescKeyMCPStatusOK)
		if f.IsEmpty {
			status = desc.TextDesc(text.TextDescKeyMCPStatusEmpty)
		}
		_, _ = fmt.Fprintf(
			&sb, desc.TextDesc(text.TextDescKeyMCPStatusFileFormat),
			f.Name, f.Tokens, status,
		)
	}

	return sb.String(), nil
}

// Add adds an entry to a context file.
//
// Parameters:
//   - entryType: the type of entry (task, decision, learning, convention)
//   - content: the entry content
//   - opts: optional fields for the entry
//
// Returns:
//   - string: confirmation message with entry type and target file
//   - error: boundary, validation, or write error
func (h *Handler) Add(
	entryType, content string, opts EntryOpts,
) (string, error) {
	if boundaryErr := validation.ValidateBoundary(
		h.ContextDir,
	); boundaryErr != nil {
		return "", boundaryErr
	}

	fileName, writeErr := entry.ValidateAndWrite(entry.Params{
		Type:        entryType,
		Content:     content,
		ContextDir:  h.ContextDir,
		Priority:    opts.Priority,
		Context:     opts.Context,
		Rationale:   opts.Rationale,
		Consequence: opts.Consequence,
		Lesson:      opts.Lesson,
		Application: opts.Application,
	})
	if writeErr != nil {
		return "", writeErr
	}

	return fmt.Sprintf(
		desc.TextDesc(text.TextDescKeyMCPAddedFormat),
		entryType, fileName,
	), nil
}

// Complete marks a task as done by number or text match.
//
// Parameters:
//   - query: task number or text fragment to match
//
// Returns:
//   - string: confirmation message with completed task text
//   - error: boundary or completion error
func (h *Handler) Complete(query string) (string, error) {
	if boundaryErr := validation.ValidateBoundary(
		h.ContextDir,
	); boundaryErr != nil {
		return "", boundaryErr
	}

	completedTask, completeErr := taskcomplete.CompleteTask(
		query, h.ContextDir,
	)
	if completeErr != nil {
		return "", completeErr
	}

	return fmt.Sprintf(
		desc.TextDesc(text.TextDescKeyMCPCompletedFormat),
		completedTask,
	), nil
}

// Drift runs drift detection and returns the report.
//
// Returns:
//   - string: formatted drift report with violations, warnings, passed
//   - error: context load error
func (h *Handler) Drift() (string, error) {
	ctx, loadErr := load.Do(h.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	report := drift.Detect(ctx)

	var sb strings.Builder
	_, _ = fmt.Fprintf(
		&sb,
		desc.TextDesc(text.TextDescKeyMCPDriftStatusFormat),
		report.Status(),
	)

	if len(report.Violations) > 0 {
		sb.WriteString(desc.TextDesc(text.TextDescKeyMCPDriftViolations))
		for _, v := range report.Violations {
			_, _ = fmt.Fprintf(
				&sb, desc.TextDesc(text.TextDescKeyMCPDriftIssueFormat),
				v.Type, v.File, v.Message,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Warnings) > 0 {
		sb.WriteString(desc.TextDesc(text.TextDescKeyMCPDriftWarnings))
		for _, w := range report.Warnings {
			_, _ = fmt.Fprintf(
				&sb, desc.TextDesc(text.TextDescKeyMCPDriftIssueFormat),
				w.Type, w.File, w.Message,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Passed) > 0 {
		sb.WriteString(desc.TextDesc(text.TextDescKeyMCPDriftPassed))
		for _, p := range report.Passed {
			_, _ = fmt.Fprintf(
				&sb, desc.TextDesc(text.TextDescKeyMCPDriftPassedFormat), p,
			)
		}
	}

	return sb.String(), nil
}

// Recall queries recent session history.
//
// Parameters:
//   - limit: max sessions to return
//   - since: only return sessions after this time (zero value = no filter)
//
// Returns:
//   - string: formatted session list with dates, projects, durations
//   - error: session discovery error
func (h *Handler) Recall(limit int, since time.Time) (string, error) {
	sessions, findErr := parser.FindSessions()
	if findErr != nil {
		return "", findErr
	}

	// Apply since filter.
	if !since.IsZero() {
		var filtered []*parser.Session
		for _, sess := range sessions {
			if sess.StartTime.After(since) ||
				sess.StartTime.Equal(since) {
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
		return desc.TextDesc(text.TextDescKeyMCPNoSessions), nil
	}

	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb,
		desc.TextDesc(text.TextDescKeyMCPSessionsFoundFormat),
		len(sessions),
	)

	for i, sess := range sessions {
		duration := sess.Duration.Round(time.Second)
		_, _ = fmt.Fprintf(
			&sb,
			desc.TextDesc(text.TextDescKeyMCPRecallItemFormat),
			i+1, sess.StartTime.Format(timeCfg.DateTimeFormat),
		)
		if sess.Project != "" {
			_, _ = fmt.Fprintf(
				&sb, desc.TextDesc(text.TextDescKeyMCPRecallProjectFormat),
				sess.Project,
			)
		}
		_, _ = fmt.Fprintf(
			&sb, desc.TextDesc(text.TextDescKeyMCPRecallDurationFormat),
			duration, sess.TurnCount,
		)
		sb.WriteString(token.NewlineLF)

		if sess.FirstUserMsg != "" {
			_, _ = fmt.Fprintf(
				&sb, desc.TextDesc(text.TextDescKeyMCPRecallFirstMsgFormat),
				sess.FirstUserMsg,
			)
			sb.WriteString(token.NewlineLF)
		}
	}

	return sb.String(), nil
}

// WatchUpdate applies a structured context-update to .context/ files.
//
// Parameters:
//   - entryType: the type of entry
//   - content: the entry content
//   - opts: optional fields for the entry
//
// Returns:
//   - string: confirmation with file name and review status
//   - error: boundary, validation, or write error
func (h *Handler) WatchUpdate(
	entryType, content string, opts EntryOpts,
) (string, error) {
	if boundaryErr := validation.ValidateBoundary(h.ContextDir); boundaryErr != nil {
		return "", boundaryErr
	}

	// Handle the "complete" type as a special case.
	if entryType == entryCfg.Complete {
		completedTask, completeErr := taskcomplete.CompleteTask(
			content, h.ContextDir)
		if completeErr != nil {
			return "", completeErr
		}
		h.Session.QueuePendingUpdate(session.PendingUpdate{
			Type:     entryType,
			Content:  content,
			QueuedAt: time.Now(),
		})
		return fmt.Sprintf(
			desc.TextDesc(text.TextDescKeyMCPWatchCompletedFormat),
			completedTask,
		) + token.NewlineLF +
			desc.TextDesc(text.TextDescKeyMCPReviewStatus), nil
	}

	fileName, writeErr := entry.ValidateAndWrite(entry.Params{
		Type:        entryType,
		Content:     content,
		ContextDir:  h.ContextDir,
		Priority:    opts.Priority,
		Context:     opts.Context,
		Rationale:   opts.Rationale,
		Consequence: opts.Consequence,
		Lesson:      opts.Lesson,
		Application: opts.Application,
	})
	if writeErr != nil {
		return "", writeErr
	}

	h.Session.RecordAdd(entryType)
	h.Session.QueuePendingUpdate(session.PendingUpdate{
		Type:    entryType,
		Content: content,
		Attrs: map[string]string{
			field.AttrFile: fileName,
		},
		QueuedAt: time.Now(),
	})

	return fmt.Sprintf(
		desc.TextDesc(text.TextDescKeyMCPWroteFormat),
		entryType, fileName,
	) + token.NewlineLF +
		desc.TextDesc(text.TextDescKeyMCPReviewStatus), nil
}

// Compact moves completed tasks to the archive section.
//
// Parameters:
//   - archive: whether to write archivable blocks to the archive file
//
// Returns:
//   - string: summary of moved tasks and cleaned sections
//   - error: boundary, context load, or write error
func (h *Handler) Compact(archive bool) (string, error) {
	if boundaryErr := validation.ValidateBoundary(
		h.ContextDir,
	); boundaryErr != nil {
		return "", boundaryErr
	}

	ctx, loadErr := load.Do(h.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	result := tidy.CompactContext(ctx)

	// Write TASKS.md changes.
	if result.TasksFileUpdate != nil {
		if writeErr := os.WriteFile(
			result.TasksFileUpdate.Path,
			result.TasksFileUpdate.Content,
			configfs.PermFile,
		); writeErr != nil {
			return "", writeErr
		}
	}

	// Write section-cleaned files.
	for _, fu := range result.SectionFileUpdates {
		if writeErr := os.WriteFile(
			fu.Path, fu.Content, configfs.PermFile,
		); writeErr != nil {
			return "", writeErr
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
			desc.TextDesc(text.DescKeyHeadingArchivedTasks),
			archiveContent,
		); archiveErr != nil {
			_, _ = fmt.Fprintf(
				&sb,
				desc.TextDesc(text.TextDescKeyMCPCompactArchiveWarning)+
					token.NewlineLF,
				archiveErr,
			)
		}
	}

	// Build response text.
	for _, taskText := range result.TasksMoved {
		_, _ = fmt.Fprintf(&sb,
			desc.TextDesc(
				text.TextDescKeyMCPCompactMovedFormat)+token.NewlineLF,
			tidy.TruncateString(taskText, cfg.TruncateLen),
		)
	}
	for _, sc := range result.SectionsCleaned {
		_, _ = fmt.Fprintf(
			&sb,
			desc.TextDesc(text.TextDescKeyMCPCompactRemovedSectFmt)+
				token.NewlineLF,
			sc.Removed, sc.FileName,
		)
	}

	if result.TotalChanges() == 0 {
		return desc.TextDesc(text.TextDescKeyMCPCompactClean), nil
	}

	_, _ = fmt.Fprintf(
		&sb,
		desc.TextDesc(text.TextDescKeyMCPCompactedFormat),
		result.TotalChanges(),
	)
	sb.WriteString(desc.TextDesc(text.TextDescKeyMCPReviewStatus))

	return sb.String(), nil
}

// Next suggests the next pending task.
//
// Returns:
//   - string: next pending task or all-complete message
//   - error: context load error
func (h *Handler) Next() (string, error) {
	ctx, loadErr := load.Do(h.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return desc.TextDesc(text.TextDescKeyMCPNoTasks), nil
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result string
	task.ForEachPending(lines, func(pt task.Pending) bool {
		result = fmt.Sprintf(
			desc.TextDesc(text.TextDescKeyMCPNextTaskFormat),
			pt.Index, pt.Content,
		)
		return true // stop after first
	})

	if result != "" {
		return result, nil
	}

	return desc.TextDesc(text.TextDescKeyMCPAllTasksComplete), nil
}

// CheckTaskCompletion checks if a recent action completed any pending
// tasks.
//
// Parameters:
//   - recentAction: description of the action to match against tasks
//
// Returns:
//   - string: matching task prompt with the completion hint, or empty
//   - error: context load error
func (h *Handler) CheckTaskCompletion(recentAction string) (string, error) {
	ctx, loadErr := load.Do(h.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return "", nil
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result string
	task.ForEachPending(lines, func(pt task.Pending) bool {
		if recentAction != "" && task.ContainsOverlap(recentAction, pt.Content) {
			result = fmt.Sprintf(
				desc.TextDesc(text.TextDescKeyMCPCheckTaskFormat)+
					token.NewlineLF+
					desc.TextDesc(text.TextDescKeyMCPCheckTaskHint),
				pt.Index, pt.Content, pt.Index,
			)
			return true
		}
		return false
	})

	return result, nil
}

// SessionEvent handles session lifecycle events (start/end).
//
// Parameters:
//   - eventType: the event type (start or end)
//   - caller: optional caller identifier for start events
//
// Returns:
//   - string: session confirmation or end-of-session summary
//   - error: unknown event type error
func (h *Handler) SessionEvent(
	eventType, caller string,
) (string, error) {
	switch eventType {
	case event.Start:
		h.Session = session.NewState(h.ContextDir)
		if caller != "" {
			return fmt.Sprintf(
				desc.TextDesc(
					text.TextDescKeyMCPSessionStartedCallerFormat,
				),
				caller, h.ContextDir,
			), nil
		}
		return fmt.Sprintf(
			desc.TextDesc(text.TextDescKeyMCPSessionStartedFormat),
			h.ContextDir,
		), nil

	case event.End:
		pending := h.Session.PendingCount()
		var sb strings.Builder
		sb.WriteString(desc.TextDesc(text.TextDescKeyMCPSessionEnding))
		sb.WriteString(token.NewlineLF)

		if pending > 0 {
			_, _ = fmt.Fprintf(
				&sb,
				desc.TextDesc(text.TextDescKeyMCPPendingUpdatesFormat),
				pending,
			)
			for i, pu := range h.Session.PendingFlush {
				_, _ = fmt.Fprintf(
					&sb,
					desc.TextDesc(text.TextDescKeyMCPPendingItemFormat)+
						token.NewlineLF,
					i+1, pu.Type,
					tidy.TruncateString(pu.Content, cfg.TruncateContentLen),
				)
			}
			sb.WriteString(
				desc.TextDesc(text.TextDescKeyMCPReviewPending),
			)
		} else {
			sb.WriteString(desc.TextDesc(text.TextDescKeyMCPNoPending))
		}

		_, _ = fmt.Fprintf(&sb,
			desc.TextDesc(text.TextDescKeyMCPSessionStatsFormat),
			h.Session.ToolCalls,
			stat.TotalAdds(h.Session.AddsPerformed),
		)

		return sb.String(), nil

	default:
		return "", mcpErr.UnknownEventType(eventType)
	}
}

// Remind lists pending session-scoped reminders.
//
// Returns:
//   - string: formatted reminder list or no-reminders message
//   - error: reminder read error
func (h *Handler) Remind() (string, error) {
	reminders, readErr := remindcore.ReadReminders()
	if readErr != nil {
		return "", readErr
	}

	if len(reminders) == 0 {
		return desc.TextDesc(text.TextDescKeyMCPNoReminders), nil
	}

	today := time.Now().Format(timeCfg.DateFormat)
	var sb strings.Builder
	_, _ = fmt.Fprintf(
		&sb,
		desc.TextDesc(text.TextDescKeyMCPRemindersFormat),
		len(reminders),
	)

	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf(
					desc.TextDesc(
						text.TextDescKeyMCPReminderNotDueFormat,
					), *r.After,
				)
			}
		}
		_, _ = fmt.Fprintf(&sb, desc.TextDesc(
			text.TextDescKeyMCPReminderItemFormat)+token.NewlineLF,
			r.ID, r.Message, annotation)
	}

	return sb.String(), nil
}
