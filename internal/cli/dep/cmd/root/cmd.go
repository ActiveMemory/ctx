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
		Use:   cmd.UseDep,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, format, external, projType)
		},
	}

	c.Flags().StringVar(
		&format,
		cFlag.Format, fmt.FormatMermaid, desc.Flag(flag.DescKeyDepsFormat),
	)
	c.Flags().BoolVar(
		&external,
		cFlag.External, false, desc.Flag(flag.DescKeyDepsExternal),
	)
	c.Flags().StringVar(
		&projType, cFlag.Type, "", desc.Flag(flag.DescKeyDepsType),
	)

	return c
}
