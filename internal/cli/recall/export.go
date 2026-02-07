//   /    Context:                     https://ctx.ist
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

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/recall/parser"
)

// maxMessagesPerPart is the maximum number of messages per exported file.
// Sessions with more messages are split into multiple parts for browser performance.
const maxMessagesPerPart = 200

// recallExportCmd returns the recall export subcommand.
//
// Returns:
//   - *cobra.Command: Command for exporting sessions to journal files
func recallExportCmd() *cobra.Command {
	var (
		all          bool
		allProjects  bool
		force        bool
		skipExisting bool
	)

	cmd := &cobra.Command{
		Use:   "export [session-id]",
		Short: "Export sessions to editable journal files",
		Long: `Export AI sessions to .context/journal/ as editable Markdown files.

Exported files include session metadata, tool usage summary, and the full
conversation. You can edit these files to add notes, highlight key moments,
or clean up the transcript.

By default, only sessions from the current project are exported. Use
--all-projects to include sessions from all projects.

By default, existing files are updated: YAML frontmatter from enrichment is
preserved, conversation content is regenerated. Use --skip-existing to leave
existing files untouched, or --force to overwrite completely.

Examples:
  ctx recall export abc123              # Export one session
  ctx recall export --all               # Export/update all sessions
  ctx recall export --all --skip-existing # Skip files that already exist
  ctx recall export --all --force       # Overwrite existing exports completely`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallExport(cmd, args, all, allProjects, force, skipExisting)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Export all sessions from current project")
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, "Include sessions from all projects")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files completely (discard frontmatter)")
	cmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip files that already exist")

	return cmd
}

// runRecallExport handles the recall export command.
func runRecallExport(cmd *cobra.Command, args []string, all, allProjects, force, skipExisting bool) error {
	if len(args) > 0 && all {
		return fmt.Errorf("cannot use --all with a session ID; use one or the other")
	}
	if len(args) == 0 && !all {
		return fmt.Errorf("please provide a session ID or use --all")
	}

	// Find sessions - filter by current project unless --all-projects is set
	var sessions []*parser.Session
	var err error
	if allProjects {
		sessions, err = parser.FindSessions()
	} else {
		cwd, cwdErr := os.Getwd()
		if cwdErr != nil {
			return fmt.Errorf("failed to get working directory: %w", cwdErr)
		}
		sessions, err = parser.FindSessionsForCWD(cwd)
	}
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
	journalDir := filepath.Join(rc.ContextDir(), "journal")
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
			if !isEmptyMessage(msg) {
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
				existing, readErr := os.ReadFile(path)
				if readErr == nil {
					if fm := extractFrontmatter(string(existing)); fm != "" {
						content = fm + "\n" + content
					}
				}
				updated++
			} else if fileExists {
				exported++
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

// extractFrontmatter returns the YAML frontmatter block from content,
// including the --- delimiters and trailing newline. Returns empty string
// if no frontmatter found.
func extractFrontmatter(content string) string {
	if !strings.HasPrefix(content, "---\n") {
		return ""
	}
	end := strings.Index(content[4:], "\n---\n")
	if end < 0 {
		return ""
	}
	// Include the full \n---\n (5 chars)
	return content[:4+end+5]
}

// formatJournalFilename generates the filename for a journal entry.
//
// Format: YYYY-MM-DD-slug-shortid.md
// Uses local time for the date.
//
// Parameters:
//   - s: Session to generate filename for
//
// Returns:
//   - string: Filename like "2026-01-15-gleaming-wobbling-sutherland-abc12345.md"
func formatJournalFilename(s *parser.Session) string {
	date := s.StartTime.Local().Format("2006-01-02")
	shortID := s.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	return fmt.Sprintf("%s-%s-%s.md", date, s.Slug, shortID)
}

// isEmptyMessage returns true if a message has no meaningful content.
func isEmptyMessage(msg parser.Message) bool {
	return msg.Text == "" && len(msg.ToolUses) == 0 && len(msg.ToolResults) == 0
}

// formatJournalEntryPart generates Markdown content for a part of a journal entry.
//
// Includes metadata, tool usage summary (on part 1 only), navigation links,
// and the conversation subset for this part.
//
// Parameters:
//   - s: Session to format
//   - messages: Subset of messages for this part
//   - startMsgIdx: Starting message index (for numbering)
//   - part: Current part number (1-indexed)
//   - totalParts: Total number of parts
//   - baseName: Base filename without extension (for navigation links)
//
// Returns:
//   - string: Markdown content for this part
func formatJournalEntryPart(
	s *parser.Session,
	messages []parser.Message,
	startMsgIdx, part, totalParts int,
	baseName string,
) string {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	// Header
	if s.Slug != "" {
		sb.WriteString(fmt.Sprintf("# %s"+nl+nl, s.Slug))
	} else {
		sb.WriteString(fmt.Sprintf("# %s"+nl+nl, baseName))
	}

	// Navigation header for multi-part sessions
	if totalParts > 1 {
		sb.WriteString(formatPartNavigation(part, totalParts, baseName, nl))
		sb.WriteString(nl + sep + nl + nl)
	}

	// Metadata (use local time) - only on part 1
	if part == 1 {
		localStart := s.StartTime.Local()
		sb.WriteString(fmt.Sprintf("**ID**: %s"+nl, s.ID))
		sb.WriteString(fmt.Sprintf("**Date**: %s"+nl, localStart.Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("**Time**: %s"+nl, localStart.Format("15:04:05")))
		sb.WriteString(fmt.Sprintf("**Duration**: %s"+nl, formatDuration(s.Duration)))
		sb.WriteString(fmt.Sprintf("**Tool**: %s"+nl, s.Tool))
		sb.WriteString(fmt.Sprintf("**Project**: %s"+nl, s.Project))
		if s.GitBranch != "" {
			sb.WriteString(fmt.Sprintf("**Branch**: %s"+nl, s.GitBranch))
		}
		if s.Model != "" {
			sb.WriteString(fmt.Sprintf("**Model**: %s"+nl, s.Model))
		}
		sb.WriteString(nl)

		// Token stats
		sb.WriteString(fmt.Sprintf("**Turns**: %d"+nl, s.TurnCount))
		sb.WriteString(fmt.Sprintf("**Tokens**: %s (in: %s, out: %s)"+nl,
			formatTokens(s.TotalTokens),
			formatTokens(s.TotalTokensIn),
			formatTokens(s.TotalTokensOut)))
		if totalParts > 1 {
			sb.WriteString(fmt.Sprintf("**Parts**: %d"+nl, totalParts))
		}
		sb.WriteString(nl + sep + nl + nl)

		// Summary section (placeholder for user to fill in)
		sb.WriteString("## Summary" + nl + nl)
		sb.WriteString("[Add your summary of this session]" + nl + nl)
		sb.WriteString(sep + nl + nl)

		// Tool usage summary
		tools := s.AllToolUses()
		if len(tools) > 0 {
			sb.WriteString("## Tool Usage" + nl + nl)
			toolCounts := make(map[string]int)
			for _, t := range tools {
				toolCounts[t.Name]++
			}
			for name, count := range toolCounts {
				sb.WriteString(fmt.Sprintf("- %s: %d"+nl, name, count))
			}
			sb.WriteString(nl + sep + nl + nl)
		}
	}

	// Conversation section
	if part == 1 {
		sb.WriteString("## Conversation" + nl + nl)
	} else {
		sb.WriteString(fmt.Sprintf("## Conversation (continued from part %d)"+nl+nl, part-1))
	}

	for i, msg := range messages {
		msgNum := startMsgIdx + i + 1
		role := "User"
		if msg.IsAssistant() {
			role = "Assistant"
		} else if len(msg.ToolResults) > 0 && msg.Text == "" {
			role = "Tool Output"
		}

		localTime := msg.Timestamp.Local()
		sb.WriteString(fmt.Sprintf("### %d. %s (%s)"+nl+nl,
			msgNum, role, localTime.Format("15:04:05")))

		if msg.Text != "" {
			text := msg.Text
			// Normalize code fences in user messages (users often type "text: ```code")
			if !msg.IsAssistant() {
				text = normalizeCodeFences(text)
			}
			sb.WriteString(text + nl + nl)
		}

		// Tool uses
		for _, t := range msg.ToolUses {
			sb.WriteString(fmt.Sprintf("🔧 **%s**"+nl, formatToolUse(t)))
		}

		// Tool results
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				sb.WriteString("❌ Error" + nl)
			}
			if tr.Content != "" {
				content := stripLineNumbers(tr.Content)
				content, reminders := extractSystemReminders(content)
				fence := fenceForContent(content)
				lines := strings.Count(content, "\n")

				if lines > 10 {
					summary := fmt.Sprintf("%d lines", lines)
					sb.WriteString(fmt.Sprintf("<details>"+nl+"<summary>%s</summary>"+nl+nl, summary))
					sb.WriteString(fmt.Sprintf("%s"+nl+"%s"+nl+"%s"+nl, fence, content, fence))
					sb.WriteString("</details>" + nl)
				} else {
					sb.WriteString(fmt.Sprintf("%s"+nl+"%s"+nl+"%s"+nl, fence, content, fence))
				}

				// Render system reminders as markdown outside the code fence
				for _, reminder := range reminders {
					sb.WriteString(fmt.Sprintf(nl+"**System Reminder**: %s"+nl, reminder))
				}
			}
		}

		if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
			sb.WriteString(nl)
		}
	}

	// Navigation footer for multi-part sessions
	if totalParts > 1 {
		sb.WriteString(nl + sep + nl + nl)
		sb.WriteString(formatPartNavigation(part, totalParts, baseName, nl))
	}

	return sb.String()
}

// formatPartNavigation generates navigation links for multi-part sessions.
func formatPartNavigation(part, totalParts int, baseName, nl string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("**Part %d of %d**", part, totalParts))

	if part > 1 || part < totalParts {
		sb.WriteString(" | ")
	}

	// Previous link
	if part > 1 {
		prevFile := baseName + ".md"
		if part > 2 {
			prevFile = fmt.Sprintf("%s-p%d.md", baseName, part-1)
		}
		sb.WriteString(fmt.Sprintf("[← Previous](%s)", prevFile))
	}

	// Separator between prev and next
	if part > 1 && part < totalParts {
		sb.WriteString(" | ")
	}

	// Next link
	if part < totalParts {
		nextFile := fmt.Sprintf("%s-p%d.md", baseName, part+1)
		sb.WriteString(fmt.Sprintf("[Next →](%s)", nextFile))
	}

	sb.WriteString(nl)
	return sb.String()
}

// fenceForContent returns the appropriate code fence for content.
//
// Uses longer fences when content contains backticks to avoid
// nested markdown rendering issues. Starts with ``` and adds
// more backticks as needed.
//
// Parameters:
//   - content: The content to be fenced
//
// Returns:
//   - string: A fence string (e.g., "```", "````", "```“")
func fenceForContent(content string) string {
	fence := "```"
	for strings.Contains(content, fence) {
		fence += "`"
	}
	return fence
}
