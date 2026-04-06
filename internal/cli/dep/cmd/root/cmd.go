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
	"github.com/ActiveMemory/ctx/internal/config/fmt"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the dep command.
//
// Flags:
//   - --format: Output format (mermaid, table, json)
//   - --external: Include external module dependencies
//   - --type: Force project type override
//
// Returns:
//   - *cobra.Command: Configured deps command with flags registered
func Cmd() *cobra.Command {
	var (
		format   string
		external bool
		projType string
	)

	short, long := desc.Command(cmd.DescKeyDep)
	c := &cobra.Command{
		Use:     cmd.UseDep,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyDep),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, format, external, projType)
		},
	}

	flagbind.StringFlagDefault(
		c, &format,
		cFlag.Format, fmt.FormatMermaid,
		flag.DescKeyDepsFormat,
	)
	flagbind.BoolFlag(
		c, &external,
		cFlag.External, flag.DescKeyDepsExternal,
	)
	flagbind.StringFlag(
		c, &projType,
		cFlag.Type, flag.DescKeyDepsType,
	)

	return c
}
