//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package amend

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the dream amend subcommand.
//
// Returns:
//   - *cobra.Command: configured amend subcommand
func Cmd() *cobra.Command {
	var (
		action string
		note   string
	)

	short, long := desc.Command(cmd.DescKeyDreamAmend)
	c := &cobra.Command{
		Use:     cmd.UseDreamAmend,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyDreamAmend),
		Args:    cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0], action, note)
		},
	}

	flagbind.StringFlag(c, &action, cFlag.Action, flag.DescKeyDreamAction)
	flagbind.StringFlag(c, &note, cFlag.Note, flag.DescKeyDreamNote)
	return c
}
