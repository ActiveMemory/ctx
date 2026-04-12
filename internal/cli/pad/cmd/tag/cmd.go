//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the pad tag subcommand.
//
// Lists all tags found across scratchpad entries, sorted
// alphabetically with occurrence counts.
//
// Returns:
//   - *cobra.Command: Configured tag subcommand
func Cmd() *cobra.Command {
	var jsonOut bool

	short, long := desc.Command(cmd.DescKeyPadTag)
	c := &cobra.Command{
		Use:     cmd.UsePadTag,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPadTag),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, jsonOut)
		},
	}

	flagbind.BoolFlag(c, &jsonOut,
		cFlag.JSON, embedFlag.DescKeyPadTagJSON,
	)

	return c
}
