//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// defaultDepth is the shallowest-only heading depth: L2 (`##`) headings only.
const defaultDepth = 2

// Cmd returns the "ctx index" command.
//
// It takes one positional argument (the Markdown file to project) and prints
// its ATX headings in file order. Deeper levels are opt-in via --depth; --json
// switches to machine-readable output.
//
// Flags:
//   - --depth: Max heading level to include (2 = ## only; 3 adds ###)
//   - --json: Emit a JSON array of {level, text} instead of lines
//
// Returns:
//   - *cobra.Command: Configured index command with flags registered.
func Cmd() *cobra.Command {
	var (
		depth      int
		jsonOutput bool
	)

	short, long := desc.Command(cmd.DescKeyIndex)
	c := &cobra.Command{
		Use:     cmd.UseIndex,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyIndex),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], depth, jsonOutput)
		},
	}

	flagbind.IntFlag(c, &depth,
		cFlag.Depth, defaultDepth, flag.DescKeyIndexDepth,
	)
	flagbind.BoolFlag(c, &jsonOutput,
		cFlag.JSON, flag.DescKeyIndexJSON,
	)

	return c
}
