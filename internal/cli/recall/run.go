//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/parse"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
	"github.com/ActiveMemory/ctx/internal/write"
)

// executeExport writes files according to the plan.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - plan: the export plan with file actions.
//   - jstate: journal state to update as files are exported.
//   - opts: export flag values.
//
// Returns:
//   - exported: number of new files written.
//   - updated: number of existing files updated (frontmatter preserved).
//   - skipped: number of files skipped (existing or locked).
func executeExport(
	cmd *cobra.Command,
	plan exportPlan,
	jstate *state.JournalState,
	opts exportOpts,
) (exported, updated, skipped int) {
	for _, fa := range plan.actions {
		if fa.action == actionLocked {
			skipped++
			write.SkipFile(cmd, fa.filename, config.FrontmatterLocked)
			continue
		}
		if fa.action == actionSkip {
			skipped++
			write.SkipFile(cmd, fa.filename, config.ReasonExists)
			continue
		}

		// Generate content, sanitizing any invalid UTF-8.
		content := strings.ToValidUTF8(
			formatJournalEntryPart(
				fa.session, fa.messages[fa.startIdx:fa.endIdx],
				fa.startIdx, fa.part, fa.totalParts, fa.baseName, fa.title,
			),
			config.Ellipsis,
		)

		fileExists := fa.action == actionRegenerate

		// Preserve enriched YAML frontmatter from the existing file.
		discard := opts.discardFrontmatter()
		if fileExists && !discard {
			existing, readErr := os.ReadFile(filepath.Clean(fa.path))
			if readErr == nil {
				if fm := extractFrontmatter(string(existing)); fm != "" {
					content = fm + config.NewlineLF + stripFrontmatter(content)
				}
			}
		}
		if fileExists && discard {
			jstate.ClearEnriched(fa.filename)
		}
		if fileExists && !discard {
			updated++
		} else {
			exported++
		}

		// Write file.
		if writeErr := os.WriteFile(
			fa.path, []byte(content), config.PermFile,
		); writeErr != nil {
			write.WarnFileErr(cmd, fa.filename, writeErr)
			continue
		}

		jstate.MarkExported(fa.filename)

		if fileExists && !discard {
			write.ExportedFile(cmd, fa.filename, config.ReasonUpdated)
		} else {
			write.ExportedFile(cmd, fa.filename, "")
		}
	}

	return exported, updated, skipped
}

// runRecallExport handles the recall export command.
//
// Parameters:
//   - cmd: Cobra command for output.
//   - args: positional arguments (optional session ID).
//   - opts: export flag values.
//
// Returns:
//   - error: non-nil on validation, scan, or write failures.
func runRecallExport(cmd *cobra.Command, args []string, opts exportOpts) error {
	// --keep-frontmatter=false implies --regenerate
	// (can't discard without regenerating).
	if !opts.keepFrontmatter {
		opts.regenerate = true
	}

	// 1. Validate flags.
	if validateErr := validateExportFlags(args, opts); validateErr != nil {
		return validateErr
	}

	// 2. Bare export (no args, no --all) → show help (T2.8).
	if len(args) == 0 && !opts.all {
		return cmd.Help()
	}

	// 3. Resolve sessions.
	sessions, scanErr := findSessions(opts.allProjects)
	if scanErr != nil {
		return ctxerr.FindSessions(scanErr)
	}

	if len(sessions) == 0 {
		write.NoSessionsForProject(cmd, opts.allProjects)
		return nil
	}

	var toExport []*parser.Session
	singleSession := false
	if opts.all {
		toExport = sessions
	} else {
		query := strings.ToLower(args[0])
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), query) ||
				strings.Contains(strings.ToLower(s.Slug), query) {
				toExport = append(toExport, s)
			}
		}
		if len(toExport) == 0 {
			return ctxerr.SessionNotFound(args[0])
		}
		if len(toExport) > 1 {
			lines := formatSessionMatchLines(toExport)
			write.AmbiguousSessionMatch(cmd, args[0], lines)
			return ctxerr.AmbiguousQuery()
		}
		singleSession = true
	}

	// 4. Ensure journal directory exists.
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)
	if mkErr := os.MkdirAll(journalDir, config.PermExec); mkErr != nil {
		return ctxerr.Mkdir(config.DirJournal, mkErr)
	}

	// 5. Load state + build index.
	jstate, loadErr := state.Load(journalDir)
	if loadErr != nil {
		return ctxerr.LoadJournalState(loadErr)
	}
	sessionIndex := buildSessionIndex(journalDir)

	// 6. Build the plan.
	plan := planExport(toExport, journalDir, sessionIndex, jstate, opts, singleSession)

	// 7. Execute renames.
	renamed := 0
	for _, rop := range plan.renameOps {
		renameJournalFiles(journalDir, rop.oldBase, rop.newBase, rop.numParts)
		jstate.Rename(
			rop.oldBase+config.ExtMarkdown, rop.newBase+config.ExtMarkdown,
		)
		renamed++
	}

	// 8. Dry-run → print summary and return.
	if opts.dryRun {
		write.ExportSummary(cmd, planCounts(plan), true)
		return nil
	}

	// 9. Confirmation prompt for regeneration.
	if plan.regenCount > 0 && !opts.yes && !singleSession {
		ok, promptErr := confirmExport(cmd, plan)
		if promptErr != nil {
			return promptErr
		}
		if !ok {
			write.Aborted(cmd)
			return nil
		}
	}

	// 10. Execute the export.
	exported, updated, skipped := executeExport(cmd, plan, jstate, opts)

	// 11. Persist journal state.
	if saveErr := jstate.Save(journalDir); saveErr != nil {
		write.WarnFileErr(cmd, config.FileJournalState, saveErr)
	}

	// 12. Print final summary.
	write.ExportFinalSummary(cmd, exported, updated, renamed, skipped)

	return nil
}

// runRecallList handles the recall list command.
//
// Finds all sessions, applies optional filters, and displays them in a
// formatted list with project, time, turn count, and preview.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - limit: maximum sessions to display (0 for unlimited)
//   - project: filter by project name (case-insensitive substring match)
//   - tool: filter by tool identifier (exact match)
//   - since: inclusive start date filter (YYYY-MM-DD)
//   - until: inclusive end date filter (YYYY-MM-DD)
//   - allProjects: if true, include sessions from all projects
//
// Returns:
//   - error: non-nil if date parsing or session scanning fails
func runRecallList(
	cmd *cobra.Command, limit int, project, tool,
	since, until string,
	allProjects bool,
) error {
	// Parse date filters
	sinceTime, sinceErr := parse.Date(since)
	if since != "" && sinceErr != nil {
		return ctxerr.InvalidDate(config.FlagSince, since, sinceErr)
	}
	untilTime, untilErr := parse.Date(until)
	if until != "" && untilErr != nil {
		return ctxerr.InvalidDate(config.FlagUntil, until, untilErr)
	}
	// --until is inclusive: advance to the end of the day
	if until != "" {
		untilTime = untilTime.Add(config.InclusiveUntilOffset)
	}

	sessions, scanErr := findSessions(allProjects)
	if scanErr != nil {
		return ctxerr.FindSessions(scanErr)
	}

	if len(sessions) == 0 {
		write.NoSessionsWithHint(cmd, allProjects)
		return nil
	}

	// Apply filters
	var filtered []*parser.Session
	for _, s := range sessions {
		if project != "" && !strings.Contains(
			strings.ToLower(s.Project), strings.ToLower(project),
		) {
			continue
		}
		if tool != "" && s.Tool != tool {
			continue
		}
		if since != "" && s.StartTime.Before(sinceTime) {
			continue
		}
		if until != "" && s.StartTime.After(untilTime) {
			continue
		}
		filtered = append(filtered, s)
	}

	if len(filtered) == 0 {
		write.NoFiltersMatch(cmd)
		return nil
	}

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	shown := 0
	if project != "" || tool != "" {
		shown = len(filtered)
	}
	write.SessionListHeader(cmd, len(sessions), shown)

	// Compute dynamic column widths from data.
	slugW, projW := len(config.ColSlug), len(config.ColProject)
	for _, s := range filtered {
		slug := truncate(s.Slug, config.SlugMaxLen)
		if len(slug) > slugW {
			slugW = len(slug)
		}
		if len(s.Project) > projW {
			projW = len(s.Project)
		}
	}

	// Print column header.
	rowFmt := fmt.Sprintf(config.TplRecallListRow, slugW, projW)
	write.SessionListRow(cmd, rowFmt,
		config.ColSlug, config.ColProject, config.ColDate,
		config.ColDuration, config.ColTurns, config.ColTokens)

	// Print sessions.
	for _, s := range filtered {
		slug := truncate(s.Slug, config.SlugMaxLen)
		dateStr := s.StartTime.Local().Format(config.DateTimeFormat)
		dur := formatDuration(s.Duration)
		turns := fmt.Sprintf("%d", s.TurnCount)
		tokens := ""
		if s.TotalTokens > 0 {
			tokens = formatTokens(s.TotalTokens)
		}
		write.SessionListRow(cmd, rowFmt,
			slug, s.Project, dateStr, dur, turns, tokens)
	}

	write.SessionListFooter(cmd, len(sessions) > len(filtered))

	return nil
}

// runRecallShow handles the recall show command.
//
// Displays detailed information about a session including metadata, token
// usage, tool usage summary, and optionally the full conversation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: session ID or slug to show (ignored if latest is true)
//   - latest: if true, show the most recent session
//   - full: if true, show complete conversation instead of preview
//   - allProjects: if true, search sessions from all projects
//
// Returns:
//   - error: non-nil if session not found or scanning fails
func runRecallShow(
	cmd *cobra.Command, args []string, latest, full, allProjects bool,
) error {
	sessions, scanErr := findSessions(allProjects)
	if scanErr != nil {
		return ctxerr.FindSessions(scanErr)
	}

	if len(sessions) == 0 {
		if allProjects {
			return ctxerr.NoSessionsFound("")
		}
		return ctxerr.NoSessionsFound(config.HintUseAllProjects)
	}

	var session *parser.Session

	switch {
	case latest:
		session = sessions[0]
	case len(args) == 0:
		return ctxerr.SessionIDRequired()
	default:
		query := strings.ToLower(args[0])
		var matches []*parser.Session
		for _, s := range sessions {
			if strings.HasPrefix(strings.ToLower(s.ID), query) ||
				strings.Contains(strings.ToLower(s.Slug), query) {
				matches = append(matches, s)
			}
		}
		if len(matches) == 0 {
			return ctxerr.SessionNotFound(args[0])
		}
		if len(matches) > 1 {
			lines := formatSessionMatchLines(matches)
			write.AmbiguousSessionMatchWithHint(
				cmd, args[0], lines, matches[0].ID[:config.SessionIDHintLen],
			)
			return ctxerr.AmbiguousQuery()
		}
		session = matches[0]
	}

	// Print session details
	write.SectionHeader(cmd, 1, session.Slug)

	write.SessionDetail(cmd, config.MetadataID, session.ID)
	write.SessionDetail(cmd, config.MetadataTool, session.Tool)
	write.SessionDetail(cmd, config.MetadataProject, session.Project)
	if session.GitBranch != "" {
		write.SessionDetail(cmd, config.MetadataBranch, session.GitBranch)
	}
	if session.Model != "" {
		write.SessionDetail(cmd, config.MetadataModel, session.Model)
	}
	write.BlankLine(cmd)

	write.SessionDetail(
		cmd, config.MetadataStarted,
		session.StartTime.Format(config.DateTimePreciseFormat),
	)
	write.SessionDetail(
		cmd, config.MetadataDuration, formatDuration(session.Duration),
	)
	write.SessionDetailInt(cmd, config.MetadataTurns, session.TurnCount)
	write.SessionDetailInt(cmd, config.MetadataMessages, len(session.Messages))
	write.BlankLine(cmd)

	write.SessionDetail(
		cmd, config.MetadataInputUsage, formatTokens(session.TotalTokensIn),
	)
	write.SessionDetail(
		cmd, config.MetadataOutputUsage, formatTokens(session.TotalTokensOut),
	)
	write.SessionDetail(
		cmd, config.MetadataTotal, formatTokens(session.TotalTokens),
	)
	write.BlankLine(cmd)

	// Tool usage summary
	tools := session.AllToolUses()
	if len(tools) > 0 {
		toolCounts := make(map[string]int)
		for _, t := range tools {
			toolCounts[t.Name]++
		}

		write.SectionHeader(cmd, 2, config.SectionToolUsage)
		for name, count := range toolCounts {
			write.ListItem(cmd, "%s: %d", name, count)
		}
		write.BlankLine(cmd)
	}

	// Messages
	if full {
		write.SectionHeader(cmd, 2, config.SectionConversation)

		for i, msg := range session.Messages {
			role := config.LabelRoleUser
			if msg.BelongsToAssistant() {
				role = config.LabelRoleAssistant
			} else if len(msg.ToolResults) > 0 && msg.Text == "" {
				role = config.LabelToolOutput
			}

			write.ConversationTurn(
				cmd, i+1, role, msg.Timestamp.Format(config.TimeFormat),
			)

			if msg.Text != "" {
				write.TextBlock(cmd, msg.Text)
			}

			for _, t := range msg.ToolUses {
				toolInfo := formatToolUse(t)
				write.SessionDetail(cmd, config.LabelTool, toolInfo)
			}

			for _, tr := range msg.ToolResults {
				if tr.IsError {
					write.Hint(cmd, config.LabelError)
				}
				if tr.Content != "" {
					content := stripLineNumbers(tr.Content)
					write.CodeBlock(cmd, content)
				}
			}

			if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
				write.BlankLine(cmd)
			}
		}
	} else {
		write.SectionHeader(cmd, 2, config.SectionConversationPreview)

		count := 0
		for _, msg := range session.Messages {
			if msg.BelongsToUser() && msg.Text != "" {
				count++
				if count > config.PreviewMaxTurns {
					write.MoreTurns(cmd, session.TurnCount-config.PreviewMaxTurns)
					break
				}
				text := msg.Text
				if len(text) > config.PreviewMaxTextLen {
					text = text[:config.PreviewMaxTextLen] + config.Ellipsis
				}
				write.NumberedItem(cmd, count, text)
			}
		}
		write.BlankLine(cmd)
		write.Hint(cmd, config.HintUseFullFlag)
	}

	return nil
}
