//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sessionevent

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx system session-event" subcommand.
//
// Returns:
//   - *cobra.Command: Configured session-event subcommand
func Cmd() *cobra.Command {
	var eventType string
	var caller string

	short, long := desc.Command(cmd.DescKeySystemSessionEvent)

	c := &cobra.Command{
		Use:     cmd.UseSystemSessionEvent,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemSessionEvent),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, eventType, caller)
		},
	}

	flagbind.StringFlag(c, &eventType,
		cFlag.Type, embFlag.DescKeySystemSessionEventType,
	)
	flagbind.StringFlag(c, &caller,
		cFlag.Caller, embFlag.DescKeySystemSessionEventCaller,
	)
	_ = c.MarkFlagRequired(cFlag.Type)
	_ = c.MarkFlagRequired(cFlag.Caller)

	return c
}
