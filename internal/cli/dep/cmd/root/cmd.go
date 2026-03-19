//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
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

	short, long := desc.CommandDesc(cmd.DescKeyDep)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyDep,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, format, external, projType)
		},
	}

	cmd.Flags().StringVar(&format, "format", "mermaid", desc.FlagDesc(flag.DescKeyDepsFormat))
	cmd.Flags().BoolVar(&external, "external", false, desc.FlagDesc(flag.DescKeyDepsExternal))
	cmd.Flags().StringVar(&projType, "type", "", desc.FlagDesc(flag.DescKeyDepsType))

	return cmd
}
