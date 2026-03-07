//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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

	short, long := assets.CommandDesc("recall.list")

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long:  long,
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
