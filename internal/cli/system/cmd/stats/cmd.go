//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system stats" subcommand.
//
// Returns:
//   - *cobra.Command: Configured stats subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(embed.CmdDescKeySystemStats)

	cmd := &cobra.Command{
		Use:   "stats",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().BoolP("follow", "f", false,
		assets.FlagDesc(embed.FlagDescKeySystemStatsFollow),
	)
	cmd.Flags().StringP("session", "s", "",
		assets.FlagDesc(embed.FlagDescKeySystemStatsSession),
	)
	cmd.Flags().IntP("last", "n", 20,
		assets.FlagDesc(embed.FlagDescKeySystemStatsLast),
	)
	cmd.Flags().BoolP("json", "j", false,
		assets.FlagDesc(embed.FlagDescKeySystemStatsJson),
	)

	return cmd
}
