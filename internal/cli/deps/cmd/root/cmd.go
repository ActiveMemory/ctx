//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the deps command.
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

	short, long := assets.CommandDesc(assets.CmdDescKeyDeps)
	cmd := &cobra.Command{
		Use:   "deps",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, format, external, projType)
		},
	}

	cmd.Flags().StringVar(&format, "format", "mermaid", assets.FlagDesc(assets.FlagDescKeyDepsFormat))
	cmd.Flags().BoolVar(&external, "external", false, assets.FlagDesc(assets.FlagDescKeyDepsExternal))
	cmd.Flags().StringVar(&projType, "type", "", assets.FlagDesc(assets.FlagDescKeyDepsType))

	return cmd
}
