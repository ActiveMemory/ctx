//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system events" subcommand.
//
// Returns:
//   - *cobra.Command: Configured events subcommand
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeySystemEvents)

	cmd := &cobra.Command{
		Use:   "events",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().StringP(
		"hook", "k", "", desc.FlagDesc(flag.FlagDescKeySystemEventsHook),
	)
	cmd.Flags().StringP(
		"session", "s", "", desc.FlagDesc(flag.FlagDescKeySystemEventsSession),
	)
	cmd.Flags().StringP(
		"event", "e", "", desc.FlagDesc(flag.FlagDescKeySystemEventsEvent),
	)
	cmd.Flags().IntP(
		"last", "n", 50, desc.FlagDesc(flag.FlagDescKeySystemEventsLast),
	)
	cmd.Flags().BoolP(
		"json", "j", false, desc.FlagDesc(flag.FlagDescKeySystemEventsJson),
	)
	cmd.Flags().BoolP(
		"all", "a", false, desc.FlagDesc(flag.FlagDescKeySystemEventsAll),
	)

	return cmd
}
