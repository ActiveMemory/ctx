//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system events" subcommand.
//
// Returns:
//   - *cobra.Command: Configured events subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(embed.CmdDescKeySystemEvents)

	cmd := &cobra.Command{
		Use:   "events",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().StringP(
		"hook", "k", "", assets.FlagDesc(embed.FlagDescKeySystemEventsHook),
	)
	cmd.Flags().StringP(
		"session", "s", "", assets.FlagDesc(embed.FlagDescKeySystemEventsSession),
	)
	cmd.Flags().StringP(
		"event", "e", "", assets.FlagDesc(embed.FlagDescKeySystemEventsEvent),
	)
	cmd.Flags().IntP(
		"last", "n", 50, assets.FlagDesc(embed.FlagDescKeySystemEventsLast),
	)
	cmd.Flags().BoolP(
		"json", "j", false, assets.FlagDesc(embed.FlagDescKeySystemEventsJson),
	)
	cmd.Flags().BoolP(
		"all", "a", false, assets.FlagDesc(embed.FlagDescKeySystemEventsAll),
	)

	return cmd
}
