//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
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

	short, long := desc.CommandDesc(cmd.DescKeyRecallShow)

	cmd := &cobra.Command{
		Use:   "show [session-id]",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, latest, full, allProjects)
		},
	}

	cmd.Flags().BoolVar(&latest, "latest", false,
		desc.FlagDesc(flag.DescKeyRecallShowLatest),
	)
	cmd.Flags().BoolVar(&full, "full", false,
		desc.FlagDesc(flag.DescKeyRecallShowFull),
	)
	cmd.Flags().BoolVar(&allProjects, "all-projects", false,
		desc.FlagDesc(flag.DescKeyRecallShowAllProjects),
	)

	return cmd
}
