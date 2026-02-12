//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

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

// runRecallExport handles the recall export command.
func runRecallExport(cmd *cobra.Command, args []string, all, allProjects, force, skipExisting bool) error {
	if len(args) > 0 && all {
		return fmt.Errorf("cannot use --all with a session ID; use one or the other")
	}
	if len(args) == 0 && !all {
		return fmt.Errorf("please provide a session ID or use --all")
	}

	sessions, err := findSessions(allProjects)
	if err != nil {
		return fmt.Errorf("failed to find sessions: %w", err)
	}

	if len(sessions) == 0 {
		if allProjects {
			cmd.Println("No sessions found.")
		} else {
			cmd.Println("No sessions found for this project. Use --all-projects to see all.")
		}
		return nil
	}

	// Determine which sessions to export
	var toExport []*parser.Session
	if all {
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
		if len(toExport) > 1 && !all {
			cmd.PrintErrf("Multiple sessions match '%s':\n", args[0])
			for _, m := range toExport {
				cmd.PrintErrf("  %s (%s) - %s\n",
					m.Slug, m.ID[:8], m.StartTime.Format("2006-01-02 15:04"))
			}
			return fmt.Errorf("ambiguous query, use a more specific ID")
		}
	}

	// Ensure journal directory exists
	journalDir := filepath.Join(rc.ContextDir(), config.DirJournal)
	if err := os.MkdirAll(journalDir, config.PermExec); err != nil {
		return fmt.Errorf("failed to create journal directory: %w", err)
	}

	// Export each session
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.FgHiBlack)

	var exported, updated, skipped int
	for _, s := range toExport {
		// Count non-empty messages to determine if splitting is needed
		var nonEmptyMsgs []parser.Message
		for _, msg := range s.Messages {
			if !emptyMessage(msg) {
				nonEmptyMsgs = append(nonEmptyMsgs, msg)
			}
		}

		// Calculate number of parts needed
		totalMsgs := len(nonEmptyMsgs)
		numParts := (totalMsgs + maxMessagesPerPart - 1) / maxMessagesPerPart
		if numParts < 1 {
			numParts = 1
		}

		baseFilename := formatJournalFilename(s)
		baseName := strings.TrimSuffix(baseFilename, ".md")

		// Export each part
		for part := 1; part <= numParts; part++ {
			filename := baseFilename
			if numParts > 1 && part > 1 {
				filename = fmt.Sprintf("%s-p%d.md", baseName, part)
			}
			path := filepath.Join(journalDir, filename)

			// Check if file exists
			_, statErr := os.Stat(path)
			fileExists := statErr == nil

			if fileExists && skipExisting {
				skipped++
				_, _ = dim.Fprintf(cmd.OutOrStdout(), "  skip %s (exists)\n", filename)
				continue
			}

			// Calculate message range for this part
			startIdx := (part - 1) * maxMessagesPerPart
			endIdx := startIdx + maxMessagesPerPart
			if endIdx > totalMsgs {
				endIdx = totalMsgs
			}

			// Generate content for this part
			content := formatJournalEntryPart(s, nonEmptyMsgs[startIdx:endIdx], startIdx, part, numParts, baseName)

			// If file exists and not --force, preserve YAML frontmatter
			if fileExists && !force {
				existing, readErr := os.ReadFile(filepath.Clean(path))
				if readErr == nil {
					if fm := extractFrontmatter(string(existing)); fm != "" {
						content = fm + "\n" + content
					}
				}
				updated++
			} else {
				exported++
			}

			// Write file
			if err := os.WriteFile(path, []byte(content), config.PermFile); err != nil {
				cmd.PrintErrf("  %s failed to write %s: %v\n", yellow("!"), filename, err)
				continue
			}

			if fileExists && !force {
				cmd.Printf("  %s %s (updated, frontmatter preserved)\n", green("✓"), filename)
			} else {
				cmd.Printf("  %s %s\n", green("✓"), filename)
			}
		}
	}

	cmd.Println()
	if exported > 0 {
		cmd.Printf("Exported %d new session(s) to %s\n", exported, journalDir)
	}
	if updated > 0 {
		cmd.Printf("Updated %d existing session(s) (YAML frontmatter preserved)\n", updated)
	}
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

	// Print sessions
	for i, s := range filtered {
		// Session number and slug
		_, _ = header.Fprintf(cmd.OutOrStdout(), "%2d. %s", i+1, s.Slug)
		_, _ = dim.Fprintf(cmd.OutOrStdout(), " (%s...)\n", s.ID[:8])

		// Details
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Project: %s", s.Project)
		if s.GitBranch != "" {
			_, _ = dim.Fprintf(cmd.OutOrStdout(), " (%s)", s.GitBranch)
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Time: %s", s.StartTime.Format("2006-01-02 15:04"))
		if s.Duration.Minutes() >= 1 {
			_, _ = dim.Fprintf(cmd.OutOrStdout(), " (%s)", formatDuration(s.Duration))
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "    Turns: %d", s.TurnCount)
		if s.TotalTokens > 0 {
			_, _ = dim.Fprintf(cmd.OutOrStdout(), ", Tokens: %s", formatTokens(s.TotalTokens))
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout())

		// Preview
		if s.FirstUserMsg != "" {
			preview := s.FirstUserMsg
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			_, _ = dim.Fprintf(cmd.OutOrStdout(), "    \"%s\"\n", preview)
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout())
	}

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

