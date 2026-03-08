//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system events" subcommand.
//
// Returns:
//   - *cobra.Command: Configured events subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemEvents)

	cmd := &cobra.Command{
		Use:   "events",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().StringP(
		"hook", "k", "", assets.FlagDesc(assets.FlagDescKeySystemEventsHook),
	)
	cmd.Flags().StringP(
		"session", "s", "", assets.FlagDesc(assets.FlagDescKeySystemEventsSession),
	)
	cmd.Flags().StringP(
		"event", "e", "", assets.FlagDesc(assets.FlagDescKeySystemEventsEvent),
	)
	cmd.Flags().IntP(
		"last", "n", 50, assets.FlagDesc(assets.FlagDescKeySystemEventsLast),
	)
	cmd.Flags().BoolP(
		"json", "j", false, assets.FlagDesc(assets.FlagDescKeySystemEventsJson),
	)
	cmd.Flags().BoolP(
		"all", "a", false, assets.FlagDesc(assets.FlagDescKeySystemEventsAll),
	)

	return cmd
}
