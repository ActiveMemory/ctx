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

// recallExportCmd returns the recall export subcommand.
//
// Returns:
//   - *cobra.Command: Command for exporting sessions to journal files
func recallExportCmd() *cobra.Command {
	var (
		all         bool
		allProjects bool
		force       bool
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

Existing files are skipped to preserve your edits. Use --force to overwrite.

Examples:
  ctx recall export abc123              # Export one session
  ctx recall export --all               # Export all sessions from this project
  ctx recall export --all --all-projects  # Export from all projects
  ctx recall export --all --force       # Overwrite existing exports`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallExport(cmd, args, all, allProjects, force)
		},
	}

	cmd.Flags().BoolVar(&all, "all", false, "Export all sessions from current project")
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, "Include sessions from all projects")
	cmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")

	return cmd
}

// runRecallExport handles the recall export command.
//
// Exports one or more sessions to .context/journal/ as Markdown files.
// Skips existing files unless force is true.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Session ID to export (ignored if all is true)
//   - all: If true, export all sessions
//   - allProjects: If true, include sessions from all projects
//   - force: If true, overwrite existing files
//
// Returns:
//   - error: Non-nil if export fails
func runRecallExport(cmd *cobra.Command, args []string, all, allProjects, force bool) error {
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
	if err := os.MkdirAll(journalDir, 0755); err != nil {
		return fmt.Errorf("failed to create journal directory: %w", err)
	}

	// Export each session
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.FgHiBlack)

	var exported, skipped int
	for _, s := range toExport {
		filename := formatJournalFilename(s)
		path := filepath.Join(journalDir, filename)

		// Check if file exists
		if _, err := os.Stat(path); err == nil && !force {
			skipped++
			dim.Fprintf(cmd.OutOrStdout(), "  skip %s (exists)\n", filename)
			continue
		}

		// Generate content
		content := formatJournalEntry(s)

		// Write file
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			cmd.PrintErrf("  %s failed to write %s: %v\n", yellow("!"), filename, err)
			continue
		}

		exported++
		cmd.Printf("  %s %s\n", green("✓"), filename)
	}

	cmd.Println()
	if exported > 0 {
		cmd.Printf("Exported %d session(s) to %s\n", exported, journalDir)
	}
	if skipped > 0 {
		dim.Fprintf(cmd.OutOrStdout(), "Skipped %d existing file(s). Use --force to overwrite.\n", skipped)
	}

	return nil
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

// formatJournalEntry generates the Markdown content for a journal entry.
//
// Includes metadata, tool usage summary, and full conversation.
//
// Parameters:
//   - s: Session to format
//
// Returns:
//   - string: Complete Markdown content
func formatJournalEntry(s *parser.Session) string {
	var sb strings.Builder
	nl := config.NewlineLF
	sep := config.Separator

	// Header
	sb.WriteString(fmt.Sprintf("# %s"+nl+nl, s.Slug))

	// Metadata (use local time)
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

	// Conversation (skip empty messages, use local time)
	sb.WriteString("## Conversation" + nl + nl)
	msgNum := 0
	for _, msg := range s.Messages {
		// Skip empty messages
		if isEmptyMessage(msg) {
			continue
		}

		msgNum++
		role := "User"
		if msg.IsAssistant() {
			role = "Assistant"
		} else if len(msg.ToolResults) > 0 && msg.Text == "" {
			// User messages with only tool results are system responses, not user input
			role = "Tool Output"
		}

		localTime := msg.Timestamp.Local()
		sb.WriteString(fmt.Sprintf("### %d. %s (%s)"+nl+nl,
			msgNum, role, localTime.Format("15:04:05")))

		if msg.Text != "" {
			sb.WriteString(msg.Text + nl + nl)
		}

		// Tool uses
		for _, t := range msg.ToolUses {
			sb.WriteString(fmt.Sprintf("🔧 **%s**"+nl, formatToolUse(t)))
		}

		// Tool results (these contain command output, file contents, etc.)
		for _, tr := range msg.ToolResults {
			if tr.IsError {
				sb.WriteString("❌ Error" + nl)
			}
			if tr.Content != "" {
				content := stripLineNumbers(tr.Content)
				sb.WriteString(fmt.Sprintf("```"+nl+"%s"+nl+"```"+nl, content))
			}
		}

		if len(msg.ToolUses) > 0 || len(msg.ToolResults) > 0 {
			sb.WriteString(nl)
		}
	}

	return sb.String()
}
