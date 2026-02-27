//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/journal/state"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// exportOpts holds all flag values for the export command.
type exportOpts struct {
	all, allProjects, force, regenerate, yes, dryRun bool
	keepFrontmatter                                  bool
}

// discardFrontmatter reports whether frontmatter should be discarded
// during regeneration, based on the combination of --keep-frontmatter
// and the deprecated --force flag.
func (o exportOpts) discardFrontmatter() bool {
	return !o.keepFrontmatter || o.force
}

// exportAction describes what will happen to a given file.
type exportAction int

const (
	actionNew        exportAction = iota // file does not exist yet
	actionRegenerate                     // file exists and will be rewritten
	actionSkip                           // file exists and will be left alone
	actionLocked                         // file is locked — never overwritten
)

// fileAction describes the planned action for a single export file (one part
// of one session).
type fileAction struct {
	session    *parser.Session
	filename   string
	path       string
	part       int
	totalParts int
	startIdx   int
	endIdx     int
	action     exportAction
	messages   []parser.Message
	slug       string
	title      string
	baseName   string
}

// exportPlan is the result of planExport: a list of per-file actions plus
// aggregate counters and any renames that need to happen first.
type exportPlan struct {
	actions     []fileAction
	newCount    int
	regenCount  int
	skipCount   int
	lockedCount int
	renameOps   []renameOp
}

// renameOp describes a dedup rename (old slug → new slug).
type renameOp struct {
	oldBase  string
	newBase  string
	numParts int
}

// findSessions returns sessions for the current project, or all projects if
// allProjects is true.
func findSessions(allProjects bool) ([]*parser.Session, error) {
	if allProjects {
		return parser.FindSessions()
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	return parser.FindSessionsForCWD(cwd)
}

// validateExportFlags checks for invalid flag combinations.
func validateExportFlags(args []string, opts exportOpts) error {
	if len(args) > 0 && opts.all {
		return fmt.Errorf("cannot use --all with a session ID; use one or the other")
	}
	if opts.regenerate && !opts.all {
		return fmt.Errorf("--regenerate requires --all (single-session export always writes)")
	}
	return nil
}

// planExport builds an exportPlan without writing any files.
func planExport(
	sessions []*parser.Session,
	journalDir string,
	sessionIndex map[string]string,
	jstate *state.JournalState,
	opts exportOpts,
	singleSession bool,
) exportPlan {
	var plan exportPlan

	for _, s := range sessions {
		// Collect non-empty messages.
		var nonEmptyMsgs []parser.Message
		for _, msg := range s.Messages {
			if !emptyMessage(msg) {
				nonEmptyMsgs = append(nonEmptyMsgs, msg)
			}
		}

		totalMsgs := len(nonEmptyMsgs)
		numParts := (totalMsgs + maxMessagesPerPart - 1) / maxMessagesPerPart
		if numParts < 1 {
			numParts = 1
		}

		// Determine title-based slug.
		var existingTitle string
		if oldFile := lookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldPath := filepath.Join(journalDir, oldFile)
			if data, readErr := os.ReadFile(filepath.Clean(oldPath)); readErr == nil {
				existingTitle = extractFrontmatterField(string(data), "title")
			}
		}
		slug, title := titleSlug(s, existingTitle)

		baseFilename := formatJournalFilename(s, slug)
		baseName := strings.TrimSuffix(baseFilename, ".md")

		// Detect renames (dedup: old slug → new slug).
		if oldFile := lookupSessionFile(sessionIndex, s.ID); oldFile != "" {
			oldBase := strings.TrimSuffix(oldFile, config.ExtMarkdown)
			if oldBase != baseName {
				plan.renameOps = append(plan.renameOps, renameOp{
					oldBase:  oldBase,
					newBase:  baseName,
					numParts: numParts,
				})
			}
		}

		// Plan each part.
		for part := 1; part <= numParts; part++ {
			filename := baseFilename
			if numParts > 1 && part > 1 {
				filename = fmt.Sprintf("%s-p%d.md", baseName, part)
			}
			path := filepath.Join(journalDir, filename)

			startIdx := (part - 1) * maxMessagesPerPart
			endIdx := startIdx + maxMessagesPerPart
			if endIdx > totalMsgs {
				endIdx = totalMsgs
			}

			_, statErr := os.Stat(path)
			fileExists := statErr == nil

			var action exportAction
			switch {
			case !fileExists:
				action = actionNew
				plan.newCount++
			case jstate.Locked(filename):
				action = actionLocked
				plan.lockedCount++
			case frontmatterHasLocked(path):
				// Frontmatter says locked — promote to state so future
				// operations skip the file without re-parsing.
				jstate.Mark(filename, "locked")
				action = actionLocked
				plan.lockedCount++
			case singleSession || opts.regenerate || opts.discardFrontmatter():
				action = actionRegenerate
				plan.regenCount++
			default:
				action = actionSkip
				plan.skipCount++
			}

			plan.actions = append(plan.actions, fileAction{
				session:    s,
				filename:   filename,
				path:       path,
				part:       part,
				totalParts: numParts,
				startIdx:   startIdx,
				endIdx:     endIdx,
				action:     action,
				messages:   nonEmptyMsgs,
				slug:       slug,
				title:      title,
				baseName:   baseName,
			})
		}
	}

	return plan
}

// printExportSummary prints what the export will (or would) do.
func printExportSummary(cmd *cobra.Command, plan exportPlan, isDryRun bool) {
	verb := "Will"
	if isDryRun {
		verb = "Would"
	}
	parts := []string{}
	if plan.newCount > 0 {
		parts = append(parts, fmt.Sprintf("export %d new", plan.newCount))
	}
	if plan.regenCount > 0 {
		parts = append(parts, fmt.Sprintf("regenerate %d existing", plan.regenCount))
	}
	if plan.skipCount > 0 {
		parts = append(parts, fmt.Sprintf("skip %d existing", plan.skipCount))
	}
	if plan.lockedCount > 0 {
		parts = append(parts, fmt.Sprintf("skip %d locked", plan.lockedCount))
	}
	if len(parts) == 0 {
		cmd.Println("Nothing to export.")
		return
	}
	cmd.Printf("%s %s.\n", verb, strings.Join(parts, ", "))
}

// confirmExport prints the plan summary and prompts for confirmation.
// Returns true if the user confirms (or if there's nothing to confirm).
func confirmExport(cmd *cobra.Command, plan exportPlan) (bool, error) {
	printExportSummary(cmd, plan, false)
	cmd.Print("Proceed? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil //nolint:goconst // trivial user input check
}

// executeExport writes files according to the plan. It returns counters for
// the final summary.
func executeExport(
	cmd *cobra.Command,
	plan exportPlan,
	jstate *state.JournalState,
	opts exportOpts,
) (exported, updated, skipped int) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.FgHiBlack)

	for _, fa := range plan.actions {
		if fa.action == actionLocked {
			skipped++
			_, _ = dim.Fprintf(cmd.OutOrStdout(),
				"  skip %s (locked)\n", fa.filename)
			continue
		}
		if fa.action == actionSkip {
			skipped++
			_, _ = dim.Fprintf(cmd.OutOrStdout(),
				"  skip %s (exists)\n", fa.filename)
			continue
		}

		// Generate content, sanitizing any invalid UTF-8.
		content := strings.ToValidUTF8(
			formatJournalEntryPart(
				fa.session, fa.messages[fa.startIdx:fa.endIdx],
				fa.startIdx, fa.part, fa.totalParts, fa.baseName, fa.title,
			),
			"...",
		)

		fileExists := fa.action == actionRegenerate

		// Preserve enriched YAML frontmatter from existing file.
		discard := opts.discardFrontmatter()
		if fileExists && !discard {
			existing, readErr := os.ReadFile(filepath.Clean(fa.path))
			if readErr == nil {
				if fm := extractFrontmatter(string(existing)); fm != "" {
					content = fm + "\n" + stripFrontmatter(content)
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
		if err := os.WriteFile(fa.path, []byte(content), config.PermFile); err != nil {
			cmd.PrintErrf("  %s failed to write %s: %v\n", yellow("!"), fa.filename, err)
			continue
		}

		jstate.MarkExported(fa.filename)

		if fileExists && !discard {
			cmd.Printf("  %s %s (updated, frontmatter preserved)\n", green("✓"), fa.filename)
		} else {
			cmd.Printf("  %s %s\n", green("✓"), fa.filename)
		}
	}

	return exported, updated, skipped
}

// runRecallExport handles the recall export command.
func runRecallExport(cmd *cobra.Command, args []string, opts exportOpts) error {
	// --keep-frontmatter=false implies --regenerate (can't discard without regenerating).
	if !opts.keepFrontmatter {
		opts.regenerate = true
	}

	// 1. Validate flags.
	if err := validateExportFlags(args, opts); err != nil {
		return err
	}

	// 2. Bare export (no args, no --all) → show help (T2.8).
	if len(args) == 0 && !opts.all {
		return cmd.Help()
	}

	// 3. Resolve sessions.
	sessions, err := findSessions(opts.allProjects)
	if err != nil {
		return fmt.Errorf("failed to find sessions: %w", err)
	}

	if len(sessions) == 0 {
		if opts.allProjects {
			cmd.Println("No sessions found.")
		} else {
			cmd.Println("No sessions found for this project. Use --all-projects to see all.")
		}
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
			return fmt.Errorf("session not found: %s", args[0])
		}
		if len(toExport) > 1 {
			cmd.PrintErrf("Multiple sessions match '%s':\n", args[0])
			for _, m := range toExport {
				cmd.PrintErrf("  %s (%s) - %s\n",
					m.Slug, m.ID[:8], m.StartTime.Format("2006-01-02 15:04"))
			}
			return fmt.Errorf("ambiguous query, use a more specific ID")
		}
		singleSession = true
	}

	// 4. Ensure journal directory exists.
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)
	if mkErr := os.MkdirAll(journalDir, config.PermExec); mkErr != nil {
		return fmt.Errorf("failed to create journal directory: %w", mkErr)
	}

	// 5. Load state + build index.
	jstate, err := state.Load(journalDir)
	if err != nil {
		return fmt.Errorf("load journal state: %w", err)
	}
	sessionIndex := buildSessionIndex(journalDir)

	// 6. Build the plan.
	plan := planExport(toExport, journalDir, sessionIndex, jstate, opts, singleSession)

	// 7. Execute renames.
	renamed := 0
	for _, rop := range plan.renameOps {
		renameJournalFiles(journalDir, rop.oldBase, rop.newBase, rop.numParts)
		jstate.Rename(rop.oldBase+config.ExtMarkdown, rop.newBase+config.ExtMarkdown)
		renamed++
	}

	// 8. Dry-run → print summary and return.
	if opts.dryRun {
		printExportSummary(cmd, plan, true)
		return nil
	}

	// 9. Confirmation prompt for regeneration.
	if plan.regenCount > 0 && !opts.yes && !singleSession {
		ok, promptErr := confirmExport(cmd, plan)
		if promptErr != nil {
			return promptErr
		}
		if !ok {
			cmd.Println("Aborted.")
			return nil
		}
	}

	// 10. Execute the export.
	exported, updated, skipped := executeExport(cmd, plan, jstate, opts)

	// 11. Persist journal state.
	if err := jstate.Save(journalDir); err != nil {
		cmd.PrintErrf("warning: failed to save journal state: %v\n", err)
	}

	// 12. Print final summary.
	cmd.Println()
	if exported > 0 {
		cmd.Printf("Exported %d new session(s) to %s\n", exported, journalDir)
	}
	if updated > 0 {
		cmd.Printf("Updated %d existing session(s) (YAML frontmatter preserved)\n", updated)
	}
	if renamed > 0 {
		cmd.Printf("Renamed %d session(s) to title-based filenames\n", renamed)
	}
	dim := color.New(color.FgHiBlack)
	if skipped > 0 {
		_, _ = dim.Fprintf(cmd.OutOrStdout(), "Skipped %d existing file(s).\n", skipped)
	}

	return nil
}

// runRecallList handles the recall list command.
//
// Finds all sessions, applies optional filters, and displays them in a
// formatted list with project, time, turn count, and preview.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - limit: Maximum sessions to display (0 for unlimited)
//   - project: Filter by project name (case-insensitive substring match)
//   - tool: Filter by tool identifier (exact match)
//   - allProjects: If true, include sessions from all projects
//
// Returns:
//   - error: Non-nil if session scanning fails
func runRecallList(cmd *cobra.Command, limit int, project, tool string, allProjects bool) error {
	sessions, err := findSessions(allProjects)
	if err != nil {
		return fmt.Errorf("failed to find sessions: %w", err)
	}

	if len(sessions) == 0 {
		if allProjects {
			cmd.Println("No sessions found.")
			cmd.Println("")
			cmd.Println("Sessions are stored in ~/.claude/projects/")
		} else {
			cmd.Println("No sessions found for this project.")
			cmd.Println("Use --all-projects to see sessions from all projects.")
		}
		return nil
	}

	// Apply filters
	var filtered []*parser.Session
	for _, s := range sessions {
		if project != "" && !strings.Contains(strings.ToLower(s.Project), strings.ToLower(project)) {
			continue
		}
		if tool != "" && s.Tool != tool {
			continue
		}
		filtered = append(filtered, s)
	}

	if len(filtered) == 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No sessions match the filters.")
		return nil
	}

	// Apply limit
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}

	// Print header
	header := color.New(color.Bold)
	dim := color.New(color.FgHiBlack)

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Found %d sessions", len(sessions))
	if project != "" || tool != "" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), " (%d shown)", len(filtered))
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	// Compute dynamic column widths from data.
	slugW, projW := len("Slug"), len("Project")
	for _, s := range filtered {
		slug := truncate(s.Slug, 36)
		if len(slug) > slugW {
			slugW = len(slug)
		}
		if len(s.Project) > projW {
			projW = len(s.Project)
		}
	}

	// Print column header.
	rowFmt := fmt.Sprintf("  %%-%ds  %%-%ds  %%-17s  %%8s  %%5s  %%7s\n", slugW, projW)
	_, _ = header.Fprintf(cmd.OutOrStdout(), rowFmt,
		"Slug", "Project", "Date", "Duration", "Turns", "Tokens")

	// Print sessions.
	for _, s := range filtered {
		slug := truncate(s.Slug, 36)
		dateStr := s.StartTime.Local().Format("2006-01-02 15:04")
		dur := formatDuration(s.Duration)
		turns := fmt.Sprintf("%d", s.TurnCount)
		tokens := ""
		if s.TotalTokens > 0 {
			tokens = formatTokens(s.TotalTokens)
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), rowFmt,
			slug, s.Project, dateStr, dur, turns, tokens)
	}

	_, _ = fmt.Fprintln(cmd.OutOrStdout())
	if len(sessions) > len(filtered) {
		_, _ = dim.Fprintf(cmd.OutOrStdout(), "Use --limit to see more sessions\n")
	}

	return nil
}

// runRecallShow handles the recall show command.
//
// Displays detailed information about a session including metadata, token
// usage, tool usage summary, and optionally the full conversation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Session ID or slug to show (ignored if latest is true)
//   - latest: If true, show the most recent session
//   - full: If true, show complete conversation instead of preview
//   - allProjects: If true, search sessions from all projects
//
// Returns:
//   - error: Non-nil if session not found or scanning fails
func runRecallShow(cmd *cobra.Command, args []string, latest, full, allProjects bool) error {
	sessions, err := findSessions(allProjects)
	if err != nil {
		return fmt.Errorf("failed to find sessions: %w", err)
	}

	if len(sessions) == 0 {
		if allProjects {
			return fmt.Errorf("no sessions found")
		}
		return fmt.Errorf("no sessions found for this project; use --all-projects to search all")
	}

	var session *parser.Session

	switch {
	case latest:
		session = sessions[0]
	case len(args) == 0:
		return fmt.Errorf("please provide a session ID or use --latest")
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
			return fmt.Errorf("session not found: %s", args[0])
		}
		if len(matches) > 1 {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Multiple sessions match '%s':\n", args[0])
			for _, m := range matches {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "  %s (%s) - %s\n",
					m.Slug, m.ID[:8], m.StartTime.Format("2006-01-02 15:04"))
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "\nUse a more specific ID (e.g., ctx recall show %s)\n", matches[0].ID[:12])
			return fmt.Errorf("ambiguous query")
		}
		session = matches[0]
	}

	// Print session details
	header := color.New(color.Bold)
	dim := color.New(color.FgHiBlack)

	_, _ = header.Fprintf(cmd.OutOrStdout(), "# %s\n", session.Slug)
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**ID**: %s\n", session.ID)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Tool**: %s\n", session.Tool)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), config.MetadataProject+" %s\n", session.Project)
	if session.GitBranch != "" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Branch**: %s\n", session.GitBranch)
	}
	if session.Model != "" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Model**: %s\n", session.Model)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Started**: %s\n", session.StartTime.Format("2006-01-02 15:04:05"))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Duration**: %s\n", formatDuration(session.Duration))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Turns**: %d\n", session.TurnCount)
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Messages**: %d\n", len(session.Messages))
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Tokens In**: %s\n", formatTokens(session.TotalTokensIn))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Tokens Out**: %s\n", formatTokens(session.TotalTokensOut))
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "**Total**: %s\n", formatTokens(session.TotalTokens))
	_, _ = fmt.Fprintln(cmd.OutOrStdout())

	// Tool usage summary
	tools := session.AllToolUses()
	if len(tools) > 0 {
		toolCounts := make(map[string]int)
		for _, t := range tools {
			toolCounts[t.Name]++
		}

		_, _ = header.Fprintf(cmd.OutOrStdout(), "## Tool Usage\n")
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
		for name, count := range toolCounts {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "- %s: %d\n", name, count)
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
	}

	// Messages
	if full {
		_, _ = header.Fprintf(cmd.OutOrStdout(), "## Conversation\n")
		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		for i, msg := range session.Messages {
			role := "User"
			roleColor := color.New(color.FgCyan, color.Bold)
			if msg.BelongsToAssistant() {
				role = "Assistant"
				roleColor = color.New(color.FgGreen, color.Bold)
			} else if len(msg.ToolResults) > 0 && msg.Text == "" {
				// User messages with only tool results are system responses
				role = "Tool Output"
				roleColor = color.New(color.FgYellow)
			}

			_, _ = roleColor.Fprintf(cmd.OutOrStdout(), "### %d. %s ", i+1, role)
			_, _ = dim.Fprintf(cmd.OutOrStdout(), "(%s)\n", msg.Timestamp.Format("15:04:05"))
			_, _ = fmt.Fprintln(cmd.OutOrStdout())

			// Show full text content - no truncation
			if msg.Text != "" {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), msg.Text)
				_, _ = fmt.Fprintln(cmd.OutOrStdout())
			}

			// Show tool uses with details
			for _, t := range msg.ToolUses {
				toolInfo := formatToolUse(t)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "🔧 **%s**\n", toolInfo)
			}

			// Show tool results
			for _, tr := range msg.ToolResults {
				if tr.IsError {
					_, _ = color.New(color.FgRed).Fprintln(cmd.OutOrStdout(), "❌ Error:")
				}
				if tr.Content != "" {
					// Strip line number prefixes and show content
					content := stripLineNumbers(tr.Content)
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), "```\n%s\n```\n", content)
				}
			}

			if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
				_, _ = fmt.Fprintln(cmd.OutOrStdout())
			}
		}
	} else {
		// Show first few user messages as preview
		_, _ = header.Fprintf(cmd.OutOrStdout(), "## Conversation Preview\n")
		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		count := 0
		for _, msg := range session.Messages {
			if msg.BelongsToUser() && msg.Text != "" {
				count++
				if count > 5 {
					_, _ = dim.Fprintf(cmd.OutOrStdout(), "... and %d more turns\n", session.TurnCount-5)
					break
				}
				text := msg.Text
				if len(text) > 100 {
					text = text[:100] + "..."
				}
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%d. %s\n", count, text)
			}
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())
		_, _ = dim.Fprintf(cmd.OutOrStdout(), "Use --full to see all messages\n")
	}

	return nil
}
