//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
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

	short, long := assets.CommandDesc(embed.CmdDescKeyRecallShow)

	cmd := &cobra.Command{
		Use:   "show [session-id]",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, latest, full, allProjects)
		},
	}

	cmd.Flags().BoolVar(&latest, "latest", false,
		assets.FlagDesc(embed.FlagDescKeyRecallShowLatest),
	)
	cmd.Flags().BoolVar(&full, "full", false,
		assets.FlagDesc(embed.FlagDescKeyRecallShowFull),
	)
	cmd.Flags().BoolVar(&allProjects, "all-projects", false,
		assets.FlagDesc(embed.FlagDescKeyRecallShowAllProjects),
	)

	return cmd
}
