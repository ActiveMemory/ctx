//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/ActiveMemory/ctx/internal/config/journal"
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

	short, long := assets.CommandDesc(embed.CmdDescKeyRecallList)

	cmd := &cobra.Command{
		Use:   "list",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, limit, project, tool, since, until, allProjects)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", journal.DefaultRecallListLimit,
		assets.FlagDesc(embed.FlagDescKeyRecallListLimit),
	)
	cmd.Flags().StringVarP(&project, "project", "p", "",
		assets.FlagDesc(embed.FlagDescKeyRecallListProject),
	)
	cmd.Flags().StringVarP(&tool, "tool", "t", "",
		assets.FlagDesc(embed.FlagDescKeyRecallListTool),
	)
	cmd.Flags().StringVar(&since, "since", "",
		assets.FlagDesc(embed.FlagDescKeyRecallListSince),
	)
	cmd.Flags().StringVar(&until, "until", "",
		assets.FlagDesc(embed.FlagDescKeyRecallListUntil),
	)
	cmd.Flags().BoolVar(&allProjects, "all-projects", false,
		assets.FlagDesc(embed.FlagDescKeyRecallListAllProjects),
	)

	return cmd
}
