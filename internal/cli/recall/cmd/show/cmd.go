//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"
)

// Cmd returns the recall show subcommand.
//
// Returns:
//   - *cobra.Command: Command for showing session details
func Cmd() *cobra.Command {
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
			return runShow(cmd, args, latest, full, allProjects)
		},
	}

	cmd.Flags().BoolVar(&latest, "latest", false, "Show the most recent session")
	cmd.Flags().BoolVar(&full, "full", false, "Show full message content")
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, "Search sessions from all projects")

	return cmd
}
