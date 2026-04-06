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
		Use:     cmd.UseSystemEvents,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemEvents),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	flagbind.StringFlagShort(c,
		cFlag.Hook, cFlag.ShortHook,
		flag.DescKeySystemEventsHook,
	)
	flagbind.StringFlagShort(c,
		cFlag.Session, cFlag.ShortSessionID,
		flag.DescKeySystemEventsSession,
	)
	flagbind.StringFlagShort(c,
		cFlag.Event, cFlag.ShortEvent,
		flag.DescKeySystemEventsEvent,
	)
	flagbind.LastJSON(c, event.DefaultLast,
		flag.DescKeySystemEventsLast,
		flag.DescKeySystemEventsJson,
	)
	flagbind.BoolFlagShort(c,
		cFlag.All, cFlag.ShortAll,
		flag.DescKeySystemEventsAll,
	)

	return c
}
