//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system stats" subcommand.
//
// Returns:
//   - *cobra.Command: Configured stats subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemStats)

	cmd := &cobra.Command{
		Use:   "stats",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().BoolP("follow", "f", false,
		assets.FlagDesc(assets.FlagDescKeySystemStatsFollow),
	)
	cmd.Flags().StringP("session", "s", "",
		assets.FlagDesc(assets.FlagDescKeySystemStatsSession),
	)
	cmd.Flags().IntP("last", "n", 20,
		assets.FlagDesc(assets.FlagDescKeySystemStatsLast),
	)
	cmd.Flags().BoolP("json", "j", false,
		assets.FlagDesc(assets.FlagDescKeySystemStatsJson),
	)

	return cmd
}
