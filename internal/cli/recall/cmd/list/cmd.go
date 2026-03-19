//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/journal"
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

	short, long := desc.CommandDesc(cmd.DescKeyRecallList)

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, limit, project, tool, since, until, allProjects)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", journal.DefaultRecallListLimit,
		desc.FlagDesc(flag.FlagDescKeyRecallListLimit),
	)
	cmd.Flags().StringVarP(&project, "project", "p", "",
		desc.FlagDesc(flag.FlagDescKeyRecallListProject),
	)
	cmd.Flags().StringVarP(&tool, "tool", "t", "",
		desc.FlagDesc(flag.FlagDescKeyRecallListTool),
	)
	cmd.Flags().StringVar(&since, "since", "",
		desc.FlagDesc(flag.FlagDescKeyRecallListSince),
	)
	cmd.Flags().StringVar(&until, "until", "",
		desc.FlagDesc(flag.FlagDescKeyRecallListUntil),
	)
	cmd.Flags().BoolVar(&allProjects, "all-projects", false,
		desc.FlagDesc(flag.FlagDescKeyRecallListAllProjects),
	)

	return cmd
}
