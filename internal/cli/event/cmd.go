//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	eventCfg "github.com/ActiveMemory/ctx/internal/config/event"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx event" top-level command.
//
// Returns:
//   - *cobra.Command: Configured event command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyEvent)

	c := &cobra.Command{
		Use:     cmd.UseEvent,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyEvent),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	flagbind.BindStringFlagShorts(c,
		[]string{cFlag.Hook, cFlag.Session, cFlag.Event},
		[]string{cFlag.ShortHook, cFlag.ShortSessionID, cFlag.ShortEvent},
		[]string{
			flag.DescKeyEventHook,
			flag.DescKeyEventSession,
			flag.DescKeyEventEvent,
		},
	)
	flagbind.LastJSON(c, eventCfg.DefaultLast,
		flag.DescKeyEventLast,
		flag.DescKeyEventJson,
	)
	flagbind.BoolFlagShort(c,
		cFlag.All, cFlag.ShortAll,
		flag.DescKeyEventAll,
	)

	return c
}
