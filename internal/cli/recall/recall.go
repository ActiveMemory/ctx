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
  export  Export sessions to editable journal files

Examples:
  ctx recall list
  ctx recall list --limit 5
  ctx recall show abc123
  ctx recall show --latest
  ctx recall export --all`,
	}

	cmd.AddCommand(recallListCmd())
	cmd.AddCommand(recallShowCmd())
	cmd.AddCommand(recallExportCmd())

	return cmd
}

// recallListCmd returns the recall list subcommand.
//
// Returns:
//   - *cobra.Command: Command for listing parsed sessions
func recallListCmd() *cobra.Command {
	var (
		limit       int
		project     string
		tool        string
		allProjects bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all parsed sessions",
		Long: `List AI sessions from the current project.

Sessions are sorted by date (newest first) and display:
  - Session slug (human-friendly name)
  - Project name
  - Start time and duration
  - Turn count (user messages)
  - Token usage

By default, only sessions from the current project are shown.
Use --all-projects to see sessions from all projects.

Examples:
  ctx recall list
  ctx recall list --limit 5
  ctx recall list --all-projects
  ctx recall list --project ctx
  ctx recall list --tool claude-code`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallList(cmd, limit, project, tool, allProjects)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum sessions to display")
	cmd.Flags().StringVarP(&project, "project", "p", "", "Filter by project name")
	cmd.Flags().StringVarP(&tool, "tool", "t", "", "Filter by tool (e.g., claude-code)")
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, "Include sessions from all projects")

	return cmd
}

// recallShowCmd returns the recall show subcommand.
//
// Returns:
//   - *cobra.Command: Command for showing session details
func recallShowCmd() *cobra.Command {
	var (
		latest      bool
		full        bool
		allProjects bool
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
By default, only searches sessions from the current project.

Examples:
  ctx recall show abc123
  ctx recall show gleaming-wobbling-sutherland
  ctx recall show --latest
  ctx recall show --latest --full
  ctx recall show abc123 --all-projects`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRecallShow(cmd, args, latest, full, allProjects)
		},
	}

	cmd.Flags().BoolVar(&latest, "latest", false, "Show the most recent session")
	cmd.Flags().BoolVar(&full, "full", false, "Show full message content")
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, "Search sessions from all projects")

	return cmd
}
