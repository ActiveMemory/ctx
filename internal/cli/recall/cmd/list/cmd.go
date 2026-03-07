//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"
)

// Cmd returns the recall list subcommand.
//
// Returns:
//   - *cobra.Command: Command for listing parsed sessions
func Cmd() *cobra.Command {
	var (
		limit       int
		project     string
		tool        string
		since       string
		until       string
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

Date filtering: --since and --until accept YYYY-MM-DD format.
Both are inclusive.

Examples:
  ctx recall list
  ctx recall list --limit 5
  ctx recall list --all-projects
  ctx recall list --project ctx
  ctx recall list --tool claude-code
  ctx recall list --since 2026-03-01
  ctx recall list --since 2026-03-01 --until 2026-03-05`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd, limit, project, tool, since, until, allProjects)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum sessions to display")
	cmd.Flags().StringVarP(&project, "project", "p", "", "Filter by project name")
	cmd.Flags().StringVarP(&tool, "tool", "t", "", "Filter by tool (e.g., claude-code)")
	cmd.Flags().StringVar(&since, "since", "", "Show sessions on or after this date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&until, "until", "", "Show sessions on or before this date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, "Include sessions from all projects")

	return cmd
}
