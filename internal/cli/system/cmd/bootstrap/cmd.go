//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system bootstrap" subcommand.
//
// Returns:
//   - *cobra.Command: Configured bootstrap subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(embed.CmdDescKeySystemBootstrap)

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().Bool("json", false,
		assets.FlagDesc(embed.FlagDescKeySystemBootstrapJson),
	)
	cmd.Flags().BoolP("quiet", "q", false,
		assets.FlagDesc(embed.FlagDescKeySystemBootstrapQuiet),
	)

	return cmd
}
