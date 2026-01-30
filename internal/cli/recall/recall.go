//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"github.com/spf13/cobra"
)

// Cmd returns the recall command with subcommands.
//
// The recall system provides commands for browsing and searching AI session
// history across multiple tools (Claude Code, Aider, etc.).
//
// Returns:
//   - *cobra.Command: The recall command with list, show, and serve subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recall",
		Short: "Browse and search AI session history",
		Long: `Browse and search AI session history from Claude Code and other tools.

The recall system parses JSONL session files and provides commands to
list sessions, view details, and search across your conversation history.

Subcommands:
  list    List all parsed sessions
  show    Show details of a specific session
  serve   Start web server for browsing (coming soon)

Examples:
  ctx recall list
  ctx recall list --limit 5
  ctx recall show abc123
  ctx recall show --latest`,
	}

	cmd.AddCommand(recallListCmd())
	cmd.AddCommand(recallShowCmd())

	return cmd
}

// recallListCmd returns the recall list subcommand.
//
// Returns:
//   - *cobra.Command: Command for listing parsed sessions
func recallListCmd() *cobra.Command {
	var (
		limit   int
		project string
		tool    string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all parsed sessions",
		Long: `List all AI sessions found in ~/.claude/projects/ and other locations.

Sessions are sorted by date (newest first) and display:
  - Session slug (human-friendly name)
  - Project name
  - Start time and duration
  - Turn count (user messages)
  - Token usage

Examples:
  ctx recall list
  ctx recall list --limit 5
  ctx recall list --project ctx
  ctx recall list --tool claude-code`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallList(cmd, limit, project, tool)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum sessions to display")
	cmd.Flags().StringVarP(&project, "project", "p", "", "Filter by project name")
	cmd.Flags().StringVarP(&tool, "tool", "t", "", "Filter by tool (e.g., claude-code)")

	return cmd
}

// recallShowCmd returns the recall show subcommand.
//
// Returns:
//   - *cobra.Command: Command for showing session details
func recallShowCmd() *cobra.Command {
	var (
		latest bool
		full   bool
	)

	cmd := &cobra.Command{
		Use:   "show [session-id]",
		Short: "Show details of a specific session",
		Long: `Show detailed information about a specific session.

The session ID can be:
  - Full session UUID
  - Partial match (first few characters)
  - Session slug name

Use --latest to show the most recent session.

Examples:
  ctx recall show abc123
  ctx recall show gleaming-wobbling-sutherland
  ctx recall show --latest
  ctx recall show --latest --full`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallShow(cmd, args, latest, full)
		},
	}

	cmd.Flags().BoolVar(&latest, "latest", false, "Show the most recent session")
	cmd.Flags().BoolVar(&full, "full", false, "Show full message content")

	return cmd
}
