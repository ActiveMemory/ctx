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

	"github.com/ActiveMemory/ctx/internal/assets"
	remindcore "github.com/ActiveMemory/ctx/internal/cli/remind/core"
	taskcomplete "github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	archiveCfg "github.com/ActiveMemory/ctx/internal/config/archive"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	entryCfg "github.com/ActiveMemory/ctx/internal/config/entry"
	configfs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/event"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	timeCfg "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entry"
	"github.com/ActiveMemory/ctx/internal/mcp/server/stat"
	"github.com/ActiveMemory/ctx/internal/mcp/session"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/tidy"
)

// Status loads context and returns a status summary.
func (h *Handler) Status() (string, error) {
	ctx, loadErr := h.loadContext()
	if loadErr != nil {
		return "", loadErr
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

	return sb.String(), nil
}

// Add adds an entry to a context file.
//
// Parameters:
//   - entryType: the type of entry (task, decision, learning, convention)
//   - content: the entry content
//   - opts: optional fields for the entry
func (h *Handler) Add(
	entryType, content string, opts EntryOpts,
) (string, error) {
	if boundaryErr := h.checkBoundary(); boundaryErr != nil {
		return "", boundaryErr
	}

	params := entry.Params{
		Type:         entryType,
		Content:      content,
		ContextDir:   h.ContextDir,
		Priority:     opts.Priority,
		Context:      opts.Context,
		Rationale:    opts.Rationale,
		Consequences: opts.Consequences,
		Lesson:       opts.Lesson,
		Application:  opts.Application,
	}

	if vErr := entry.Validate(params, nil); vErr != nil {
		return "", vErr
	}

	if wErr := entry.Write(params); wErr != nil {
		return "", wErr
	}

	fileName := entryCfg.ToCtxFile[strings.ToLower(params.Type)]
	return fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyMCPAddedFormat),
		entryType, fileName,
	), nil
}

// Complete marks a task as done by number or text match.
func (h *Handler) Complete(query string) (string, error) {
	if boundaryErr := h.checkBoundary(); boundaryErr != nil {
		return "", boundaryErr
	}

	completedTask, completeErr := taskcomplete.CompleteTask(
		query, h.ContextDir,
	)
	if completeErr != nil {
		return "", completeErr
	}

	return fmt.Sprintf(
		assets.TextDesc(assets.TextDescKeyMCPCompletedFormat),
		completedTask,
	), nil
}

// Drift runs drift detection and returns the report.
func (h *Handler) Drift() (string, error) {
	ctx, loadErr := h.loadContext()
	if loadErr != nil {
		return "", loadErr
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

	return sb.String(), nil
}

// Recall queries recent session history.
//
// Parameters:
//   - limit: max sessions to return
//   - since: only return sessions after this time (zero value = no filter)
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
		return assets.TextDesc(assets.TextDescKeyMCPNoSessions), nil
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

	return sb.String(), nil
}

// WatchUpdate applies a structured context-update to .context/ files.
//
// Parameters:
//   - entryType: the type of entry
//   - content: the entry content
//   - opts: optional fields for the entry
func (h *Handler) WatchUpdate(
	entryType, content string, opts EntryOpts,
) (string, error) {
	if boundaryErr := h.checkBoundary(); boundaryErr != nil {
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
			assets.TextDesc(assets.TextDescKeyMCPWatchCompletedFormat),
			completedTask,
		) + token.NewlineLF +
			assets.TextDesc(assets.TextDescKeyMCPReviewStatus), nil
	}

	params := entry.Params{
		Type:         entryType,
		Content:      content,
		ContextDir:   h.ContextDir,
		Priority:     opts.Priority,
		Context:      opts.Context,
		Rationale:    opts.Rationale,
		Consequences: opts.Consequences,
		Lesson:       opts.Lesson,
		Application:  opts.Application,
	}

	if vErr := entry.Validate(params, nil); vErr != nil {
		return "", vErr
	}

	if wErr := entry.Write(params); wErr != nil {
		return "", wErr
	}

	fileName := entryCfg.ToCtxFile[strings.ToLower(params.Type)]

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
		assets.TextDesc(assets.TextDescKeyMCPWroteFormat),
		entryType, fileName,
	) + token.NewlineLF +
		assets.TextDesc(assets.TextDescKeyMCPReviewStatus), nil
}

// Compact moves completed tasks to the archive section.
func (h *Handler) Compact(archive bool) (string, error) {
	if boundaryErr := h.checkBoundary(); boundaryErr != nil {
		return "", boundaryErr
	}

	ctx, loadErr := h.loadContext()
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
		return assets.TextDesc(assets.TextDescKeyMCPCompactClean), nil
	}

	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPCompactedFormat),
		result.TotalChanges(),
	)
	sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPReviewStatus))

	return sb.String(), nil
}

// Next suggests the next pending task.
func (h *Handler) Next() (string, error) {
	ctx, loadErr := h.loadContext()
	if loadErr != nil {
		return "", loadErr
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return assets.TextDesc(assets.TextDescKeyMCPNoTasks), nil
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result string
	eachPendingTask(lines, func(pt pendingTask) bool {
		result = fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPNextTaskFormat),
			pt.Index, pt.Content,
		)
		return true // stop after first
	})

	if result != "" {
		return result, nil
	}

	return assets.TextDesc(assets.TextDescKeyMCPAllTasksComplete), nil
}

// CheckTaskCompletion checks if a recent action completed any pending
// tasks.
func (h *Handler) CheckTaskCompletion(recentAction string) (string, error) {
	ctx, loadErr := h.loadContext()
	if loadErr != nil {
		return "", loadErr
	}

	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile == nil {
		return "", nil
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result string
	eachPendingTask(lines, func(pt pendingTask) bool {
		if recentAction != "" && containsOverlap(recentAction, pt.Content) {
			result = fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPCheckTaskFormat)+
					token.NewlineLF+
					assets.TextDesc(assets.TextDescKeyMCPCheckTaskHint),
				pt.Index, pt.Content, pt.Index,
			)
			return true
		}
		return false
	})

	return result, nil
}

// SessionEvent handles session lifecycle events (start/end).
func (h *Handler) SessionEvent(
	eventType, caller string,
) (string, error) {
	switch eventType {
	case event.Start:
		h.Session = session.NewState(h.ContextDir)
		if caller != "" {
			return fmt.Sprintf(
				assets.TextDesc(
					assets.TextDescKeyMCPSessionStartedCallerFormat,
				),
				caller, h.ContextDir,
			), nil
		}
		return fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyMCPSessionStartedFormat),
			h.ContextDir,
		), nil

	case event.End:
		pending := h.Session.PendingCount()
		var sb strings.Builder
		sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPSessionEnding))
		sb.WriteString(token.NewlineLF)

		if pending > 0 {
			_, _ = fmt.Fprintf(
				&sb,
				assets.TextDesc(assets.TextDescKeyMCPPendingUpdatesFormat),
				pending,
			)
			for i, pu := range h.Session.PendingFlush {
				_, _ = fmt.Fprintf(
					&sb,
					assets.TextDesc(assets.TextDescKeyMCPPendingItemFormat)+
						token.NewlineLF,
					i+1, pu.Type,
					tidy.TruncateString(pu.Content, cfg.TruncateContentLen),
				)
			}
			sb.WriteString(
				assets.TextDesc(assets.TextDescKeyMCPReviewPending),
			)
		} else {
			sb.WriteString(assets.TextDesc(assets.TextDescKeyMCPNoPending))
		}

		_, _ = fmt.Fprintf(&sb,
			assets.TextDesc(assets.TextDescKeyMCPSessionStatsFormat),
			h.Session.ToolCalls,
			stat.TotalAdds(h.Session.AddsPerformed),
		)

		return sb.String(), nil

	default:
		return "", fmt.Errorf(
			assets.TextDesc(assets.TextDescKeyMCPUnknownEventType),
			eventType,
		)
	}
}

// Remind lists pending session-scoped reminders.
func (h *Handler) Remind() (string, error) {
	reminders, readErr := remindcore.ReadReminders()
	if readErr != nil {
		return "", readErr
	}

	if len(reminders) == 0 {
		return assets.TextDesc(assets.TextDescKeyMCPNoReminders), nil
	}

	today := time.Now().Format(timeCfg.DateFormat)
	var sb strings.Builder
	_, _ = fmt.Fprintf(
		&sb,
		assets.TextDesc(assets.TextDescKeyMCPRemindersFormat),
		len(reminders),
	)

	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf(
					assets.TextDesc(
						assets.TextDescKeyMCPReminderNotDueFormat,
					), *r.After,
				)
			}
		}
		_, _ = fmt.Fprintf(&sb, assets.TextDesc(
			assets.TextDescKeyMCPReminderItemFormat)+token.NewlineLF,
			r.ID, r.Message, annotation)
	}

	return sb.String(), nil
}

// ParseRecallSince parses a since date string into a time.Time.
// Returns zero time if the string is empty.
func ParseRecallSince(since string) (time.Time, error) {
	if since == "" {
		return time.Time{}, nil
	}
	return time.Parse(timeCfg.DateFormat, since)
}
