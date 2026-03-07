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

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, assets.FlagDesc("recall.list.limit"))
	cmd.Flags().StringVarP(&project, "project", "p", "", assets.FlagDesc("recall.list.project"))
	cmd.Flags().StringVarP(&tool, "tool", "t", "", assets.FlagDesc("recall.list.tool"))
	cmd.Flags().StringVar(&since, "since", "", assets.FlagDesc("recall.list.since"))
	cmd.Flags().StringVar(&until, "until", "", assets.FlagDesc("recall.list.until"))
	cmd.Flags().BoolVar(&allProjects, "all-projects", false, assets.FlagDesc("recall.list.all-projects"))

	return cmd
}
