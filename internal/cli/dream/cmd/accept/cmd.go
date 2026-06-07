//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package accept

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the dream accept subcommand.
//
// Returns:
//   - *cobra.Command: configured accept subcommand
func Cmd() *cobra.Command {
	var note string

	short, long := desc.Command(cmd.DescKeyDreamAccept)
	c := &cobra.Command{
		Use:     cmd.UseDreamAccept,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyDreamAccept),
		Args:    cobra.ExactArgs(1),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			return Run(cobraCmd, args[0], note)
		},
	}

	flagbind.StringFlag(c, &note, cFlag.Note, flag.DescKeyDreamNote)
	return c
}
