//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package events

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/event"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx system events" subcommand.
//
// Returns:
//   - *cobra.Command: Configured events subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemEvents)

	c := &cobra.Command{
		Use:   cmd.UseSystemEvents,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	c.Flags().StringP(
		cFlag.Hook, cFlag.ShortHook, "",
		desc.Flag(flag.DescKeySystemEventsHook),
	)
	c.Flags().StringP(
		cFlag.Session, cFlag.ShortSessionID, "",
		desc.Flag(flag.DescKeySystemEventsSession),
	)
	c.Flags().StringP(
		cFlag.Event, cFlag.ShortEvent, "",
		desc.Flag(flag.DescKeySystemEventsEvent),
	)
	flagbind.LastJSON(c, event.DefaultEventsLast,
		flag.DescKeySystemEventsLast, flag.DescKeySystemEventsJson,
	)
	c.Flags().BoolP(
		cFlag.All, cFlag.ShortAll, false,
		desc.Flag(flag.DescKeySystemEventsAll),
	)

	return c
}
