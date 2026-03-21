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
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
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
		cflag.Hook, cflag.ShortHook, "", desc.Flag(flag.DescKeySystemEventsHook),
	)
	c.Flags().StringP(
		cflag.Session, cflag.ShortSessionID, "", desc.Flag(flag.DescKeySystemEventsSession),
	)
	c.Flags().StringP(
		cflag.Event, cflag.ShortEvent, "", desc.Flag(flag.DescKeySystemEventsEvent),
	)
	c.Flags().IntP(
		cflag.Last, cflag.ShortLast, 50, desc.Flag(flag.DescKeySystemEventsLast),
	)
	c.Flags().BoolP(
		cflag.JSON, cflag.ShortJSON, false, desc.Flag(flag.DescKeySystemEventsJson),
	)
	c.Flags().BoolP(
		cflag.All, cflag.ShortAll, false, desc.Flag(flag.DescKeySystemEventsAll),
	)

	return c
}
