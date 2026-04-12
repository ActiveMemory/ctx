//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	embedCmd "github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/wrap"
)

// Cmd returns the "ctx fmt" command for formatting context files.
//
// Returns:
//   - *cobra.Command: Configured fmt command with flags registered
func Cmd() *cobra.Command {
	var (
		width int
		check bool
	)

	short, long := desc.Command(embedCmd.DescKeyFmt)

	c := &cobra.Command{
		Use:     embedCmd.UseFmt,
		Short:   short,
		Long:    long,
		Example: desc.Example(embedCmd.DescKeyFmt),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, width, check)
		},
	}

	flagbind.IntFlag(
		c, &width,
		cFlag.Width, wrap.DefaultWidth, embedFlag.DescKeyFmtWidth,
	)
	flagbind.BoolFlag(
		c, &check,
		cFlag.Check, embedFlag.DescKeyFmtCheck,
	)

	return c
}
