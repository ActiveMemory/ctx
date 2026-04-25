//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	remindStore "github.com/ActiveMemory/ctx/internal/cli/remind/core/store"
	taskComplete "github.com/ActiveMemory/ctx/internal/cli/task/core/complete"
	cfgArchive "github.com/ActiveMemory/ctx/internal/config/archive"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/mcp/event"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/drift"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/entry"
	errMcp "github.com/ActiveMemory/ctx/internal/err/mcp"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/journal/parser"
	"github.com/ActiveMemory/ctx/internal/mcp/handler/task"
	"github.com/ActiveMemory/ctx/internal/mcp/server/stat"
	"github.com/ActiveMemory/ctx/internal/tidy"
)

// Status loads context and returns a status summary.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//
// Returns:
//   - string: formatted status text with file list and token counts
//   - error: context load error
func Status(d *entity.MCPDeps) (string, error) {
	ctx, loadErr := load.Do(d.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	var sb strings.Builder
	io.SafeFprintf(
		&sb,
		desc.Text(text.DescKeyMCPStatusContextFormat), ctx.Dir,
	)
	io.SafeFprintf(
		&sb,
		desc.Text(text.DescKeyMCPStatusFilesFormat), len(ctx.Files),
	)
	io.SafeFprintf(
		&sb,
		desc.Text(text.DescKeyMCPStatusUsageFormat), ctx.TotalTokens,
	)

	for _, f := range ctx.Files {
		status := desc.Text(text.DescKeyMCPStatusOK)
		if f.IsEmpty {
			status = desc.Text(text.DescKeyMCPStatusEmpty)
		}
		_, _ = fmt.Fprintf(
			&sb, desc.Text(text.DescKeyMCPStatusFileFormat),
			f.Name, f.Tokens, status,
		)
	}

	return sb.String(), nil
}

// Add adds an entry to a context file.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//   - entryType: the type of entry (task, decision, learning, convention)
//   - content: the entry content
//   - opts: optional fields for the entry
//
// Returns:
//   - string: confirmation message with entry type and target file
//   - error: validation or write error
func Add(
	d *entity.MCPDeps,
	entryType, content string, opts entity.EntryOpts,
) (string, error) {
	if writeErr := entry.ValidateAndWrite(entity.EntryParams{
		Type:        entryType,
		Content:     content,
		ContextDir:  d.ContextDir,
		Priority:    opts.Priority,
		Section:     opts.Section,
		SessionID:   opts.SessionID,
		Branch:      opts.Branch,
		Commit:      opts.Commit,
		Context:     opts.Context,
		Rationale:   opts.Rationale,
		Consequence: opts.Consequence,
		Lesson:      opts.Lesson,
		Application: opts.Application,
	}); writeErr != nil {
		return "", writeErr
	}

	return fmt.Sprintf(
		desc.Text(text.DescKeyMCPAddedFormat),
		entryType, cfgEntry.MustCtxFile(entryType),
	), nil
}

// Complete marks a task as done by number or text match.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//   - query: task number or text fragment to match
//
// Returns:
//   - string: confirmation message with completed task text
//   - error: completion error
func Complete(d *entity.MCPDeps, query string) (string, error) {
	completedTask, _, completeErr := taskComplete.Complete(
		query, d.ContextDir,
	)
	if completeErr != nil {
		return "", completeErr
	}

	return fmt.Sprintf(
		desc.Text(text.DescKeyMCPCompletedFormat),
		completedTask,
	), nil
}

// Drift runs drift detection and returns the report.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//
// Returns:
//   - string: formatted drift report with violations, warnings, passed
//   - error: context load error
func Drift(d *entity.MCPDeps) (string, error) {
	ctx, loadErr := load.Do(d.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	report := drift.Detect(ctx)

	var sb strings.Builder
	io.SafeFprintf(
		&sb,
		desc.Text(text.DescKeyMCPDriftStatusFormat),
		report.Status(),
	)

	if len(report.Violations) > 0 {
		sb.WriteString(desc.Text(text.DescKeyMCPDriftViolations))
		for _, v := range report.Violations {
			_, _ = fmt.Fprintf(
				&sb, desc.Text(text.DescKeyMCPDriftIssueFormat),
				v.Type, v.File, v.Message,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Warnings) > 0 {
		sb.WriteString(desc.Text(text.DescKeyMCPDriftWarnings))
		for _, w := range report.Warnings {
			_, _ = fmt.Fprintf(
				&sb, desc.Text(text.DescKeyMCPDriftIssueFormat),
				w.Type, w.File, w.Message,
			)
		}
		sb.WriteString(token.NewlineLF)
	}

	if len(report.Passed) > 0 {
		sb.WriteString(desc.Text(text.DescKeyMCPDriftOK))
		for _, p := range report.Passed {
			_, _ = fmt.Fprintf(
				&sb, desc.Text(text.DescKeyMCPDriftOKFormat), p,
			)
		}
	}

	return sb.String(), nil
}

// Recall queries recent session history.
//
// Parameters:
//   - d: runtime dependencies (unused, kept for signature uniformity)
//   - limit: max sessions to return
//   - since: only return sessions after this time (zero value = no filter)
//
// Returns:
//   - string: formatted session list with dates, projects, durations
//   - error: session discovery error
func Recall(
	_ *entity.MCPDeps, limit int, since time.Time,
) (string, error) {
	sessions, findErr := parser.FindSessions()
	if findErr != nil {
		return "", findErr
	}

	// Apply since filter.
	if !since.IsZero() {
		var filtered []*entity.Session
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
		return desc.Text(text.DescKeyMCPNoSessions), nil
	}

	var sb strings.Builder
	io.SafeFprintf(&sb,
		desc.Text(text.DescKeyMCPSessionsFoundFormat),
		len(sessions),
	)

	for i, sess := range sessions {
		duration := sess.Duration.Round(time.Second)
		_, _ = fmt.Fprintf(
			&sb,
			desc.Text(text.DescKeyMCPJournalSourceItemFormat),
			i+1, sess.StartTime.Format(cfgTime.DateTimeFmt),
		)
		if sess.Project != "" {
			_, _ = fmt.Fprintf(
				&sb, desc.Text(text.DescKeyMCPJournalSourceProjectFormat),
				sess.Project,
			)
		}
		_, _ = fmt.Fprintf(
			&sb, desc.Text(text.DescKeyMCPJournalSourceDurationFormat),
			duration, sess.TurnCount,
		)
		sb.WriteString(token.NewlineLF)

		if sess.FirstUserMsg != "" {
			_, _ = fmt.Fprintf(
				&sb, desc.Text(text.DescKeyMCPJournalSourceFirstMsgFormat),
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
//   - d: runtime dependencies carrying the context directory and session
//   - entryType: the type of entry
//   - content: the entry content
//   - opts: optional fields for the entry
//
// Returns:
//   - string: confirmation with file name and review status
//   - error: validation or write error
func WatchUpdate(
	d *entity.MCPDeps,
	entryType, content string, opts entity.EntryOpts,
) (string, error) {
	// Handle the "complete" type as a special case.
	if entryType == cfgEntry.Complete {
		completedTask, _, completeErr := taskComplete.Complete(
			content, d.ContextDir)
		if completeErr != nil {
			return "", completeErr
		}
		d.Session.QueuePendingUpdate(entity.PendingUpdate{
			Type:     entryType,
			Content:  content,
			QueuedAt: time.Now(),
		})
		return fmt.Sprintf(
			desc.Text(text.DescKeyMCPFormatWatchCompleted),
			completedTask,
		) + token.NewlineLF +
			desc.Text(text.DescKeyMCPReviewStatus), nil
	}

	if writeErr := entry.ValidateAndWrite(entity.EntryParams{
		Type:        entryType,
		Content:     content,
		ContextDir:  d.ContextDir,
		Priority:    opts.Priority,
		Section:     opts.Section,
		SessionID:   opts.SessionID,
		Branch:      opts.Branch,
		Commit:      opts.Commit,
		Context:     opts.Context,
		Rationale:   opts.Rationale,
		Consequence: opts.Consequence,
		Lesson:      opts.Lesson,
		Application: opts.Application,
	}); writeErr != nil {
		return "", writeErr
	}

	d.Session.RecordAdd(entryType)
	d.Session.QueuePendingUpdate(entity.PendingUpdate{
		Type:    entryType,
		Content: content,
		Attrs: map[string]string{
			field.AttrFile: cfgEntry.MustCtxFile(entryType),
		},
		QueuedAt: time.Now(),
	})

	return fmt.Sprintf(
		desc.Text(text.DescKeyMCPFormatWrote),
		entryType, cfgEntry.MustCtxFile(entryType),
	) + token.NewlineLF +
		desc.Text(text.DescKeyMCPReviewStatus), nil
}

// Compact moves completed tasks to the archive section.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//   - archive: whether to write archivable blocks to the archive file
//
// Returns:
//   - string: summary of moved tasks and cleaned sections
//   - error: context load or write error
func Compact(d *entity.MCPDeps, archive bool) (string, error) {
	ctx, loadErr := load.Do(d.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	result := tidy.CompactContext(ctx)

	// Write TASKS.md changes.
	if result.TasksFileUpdate != nil {
		if writeErr := io.SafeWriteFile(
			result.TasksFileUpdate.Path,
			result.TasksFileUpdate.Content,
			cfgFs.PermFile,
		); writeErr != nil {
			return "", writeErr
		}
	}

	// Write section-cleaned files.
	for _, fu := range result.SectionFileUpdates {
		if writeErr := io.SafeWriteFile(
			fu.Path, fu.Content, cfgFs.PermFile,
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
			cfgArchive.ScopeTasks,
			desc.Text(text.DescKeyHeadingArchivedTasks),
			archiveContent,
		); archiveErr != nil {
			_, _ = fmt.Fprintf(
				&sb,
				desc.Text(text.DescKeyMCPCompactArchiveWarning)+
					token.NewlineLF,
				archiveErr,
			)
		}
	}

	// Build response text.
	for _, taskText := range result.TasksMoved {
		io.SafeFprintf(&sb,
			desc.Text(
				text.DescKeyMCPCompactMovedFormat)+token.NewlineLF,
			tidy.TruncateString(taskText, token.TruncateLen),
		)
	}
	for _, sc := range result.SectionsCleaned {
		_, _ = fmt.Fprintf(
			&sb,
			desc.Text(text.DescKeyMCPCompactRemovedSectFmt)+
				token.NewlineLF,
			sc.Removed, sc.FileName,
		)
	}

	if result.TotalChanges() == 0 {
		return desc.Text(text.DescKeyMCPCompactClean), nil
	}

	io.SafeFprintf(
		&sb,
		desc.Text(text.DescKeyMCPFormatCompacted),
		result.TotalChanges(),
	)
	sb.WriteString(desc.Text(text.DescKeyMCPReviewStatus))

	return sb.String(), nil
}

// Next suggests the next pending task.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//
// Returns:
//   - string: next pending task or all-complete message
//   - error: context load error
func Next(d *entity.MCPDeps) (string, error) {
	ctx, loadErr := load.Do(d.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	tasksFile := ctx.File(cfgCtx.Task)
	if tasksFile == nil {
		return desc.Text(text.DescKeyMCPNoTasks), nil
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result string
	task.ForEachPending(lines, func(pt task.Pending) bool {
		result = fmt.Sprintf(
			desc.Text(text.DescKeyMCPNextTaskFormat),
			pt.Index, pt.Content,
		)
		return true // stop after first
	})

	if result != "" {
		return result, nil
	}

	return desc.Text(text.DescKeyMCPAllTasksComplete), nil
}

// CheckTaskCompletion checks if a recent action completed any pending
// tasks.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory
//   - recentAction: description of the action to match against tasks
//
// Returns:
//   - string: matching task prompt with the completion hint, or empty
//   - error: context load error
func CheckTaskCompletion(
	d *entity.MCPDeps, recentAction string,
) (string, error) {
	ctx, loadErr := load.Do(d.ContextDir)
	if loadErr != nil {
		return "", loadErr
	}

	tasksFile := ctx.File(cfgCtx.Task)
	if tasksFile == nil {
		return "", nil
	}

	lines := strings.Split(string(tasksFile.Content), token.NewlineLF)

	var result string
	task.ForEachPending(lines, func(pt task.Pending) bool {
		if recentAction != "" && task.ContainsOverlap(recentAction, pt.Content) {
			result = fmt.Sprintf(
				desc.Text(text.DescKeyMCPCheckTaskFormat)+
					token.NewlineLF+
					desc.Text(text.DescKeyMCPCheckTaskHint),
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
// On a "start" event the session state is reset to a fresh value
// and marked as started, so subsequent governance checks measure
// against the new session.
//
// Parameters:
//   - d: runtime dependencies carrying the context directory and session
//   - eventType: the event type (start or end)
//   - caller: optional caller identifier for start events
//
// Returns:
//   - string: session confirmation or end-of-session summary
//   - error: unknown event type error
func SessionEvent(
	d *entity.MCPDeps, eventType, caller string,
) (string, error) {
	switch eventType {
	case event.Start:
		d.Session = entity.NewMCPSession()
		d.Session.RecordSessionStart()
		if caller != "" {
			return fmt.Sprintf(
				desc.Text(
					text.DescKeyMCPSessionStartedCallerFormat,
				),
				caller, d.ContextDir,
			), nil
		}
		return fmt.Sprintf(
			desc.Text(text.DescKeyMCPSessionStartedFormat),
			d.ContextDir,
		), nil

	case event.End:
		pending := d.Session.PendingCount()
		var sb strings.Builder
		sb.WriteString(desc.Text(text.DescKeyMCPSessionEnding))
		sb.WriteString(token.NewlineLF)

		if pending > 0 {
			_, _ = fmt.Fprintf(
				&sb,
				desc.Text(text.DescKeyMCPPendingUpdatesFormat),
				pending,
			)
			for i, pu := range d.Session.PendingFlush {
				_, _ = fmt.Fprintf(
					&sb,
					desc.Text(text.DescKeyMCPFormatPendingItem)+
						token.NewlineLF,
					i+1, pu.Type,
					tidy.TruncateString(pu.Content, token.TruncateContentLen),
				)
			}
			sb.WriteString(
				desc.Text(text.DescKeyMCPReviewPending),
			)
		} else {
			sb.WriteString(desc.Text(text.DescKeyMCPNoPending))
		}

		io.SafeFprintf(&sb,
			desc.Text(text.DescKeyMCPFormatSessionStats),
			d.Session.ToolCalls,
			stat.TotalAdds(d.Session.AddsPerformed),
		)

		return sb.String(), nil

	default:
		return "", errMcp.UnknownEventType(eventType)
	}
}

// Remind lists pending session-scoped reminders.
//
// Parameters:
//   - d: runtime dependencies (unused, kept for signature uniformity)
//
// Returns:
//   - string: formatted reminder list or no-reminders message
//   - error: reminder read error
func Remind(_ *entity.MCPDeps) (string, error) {
	reminders, readErr := remindStore.Read()
	if readErr != nil {
		return "", readErr
	}

	if len(reminders) == 0 {
		return desc.Text(text.DescKeyMCPNoReminders), nil
	}

	today := time.Now().Format(cfgTime.DateFormat)
	var sb strings.Builder
	io.SafeFprintf(
		&sb,
		desc.Text(text.DescKeyMCPRemindersFormat),
		len(reminders),
	)

	for _, r := range reminders {
		annotation := ""
		if r.After != nil {
			if *r.After > today {
				annotation = fmt.Sprintf(
					desc.Text(
						text.DescKeyMCPFormatReminderNotDue,
					), *r.After,
				)
			}
		}
		io.SafeFprintf(&sb, desc.Text(
			text.DescKeyMCPFormatReminderItem)+token.NewlineLF,
			r.ID, r.Message, annotation)
	}

	return sb.String(), nil
}
