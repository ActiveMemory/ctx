//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system stats" subcommand.
//
// Returns:
//   - *cobra.Command: Configured stats subcommand
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeySystemStats)

	cmd := &cobra.Command{
		Use:   "stats",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().BoolP("follow", "f", false,
		desc.FlagDesc(flag.DescKeySystemStatsFollow),
	)
	cmd.Flags().StringP("session", "s", "",
		desc.FlagDesc(flag.DescKeySystemStatsSession),
	)
	cmd.Flags().IntP("last", "n", 20,
		desc.FlagDesc(flag.DescKeySystemStatsLast),
	)
	cmd.Flags().BoolP("json", "j", false,
		desc.FlagDesc(flag.DescKeySystemStatsJson),
	)

	return cmd
}
