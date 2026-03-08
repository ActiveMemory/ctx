//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system bootstrap" subcommand.
//
// Returns:
//   - *cobra.Command: Configured bootstrap subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemBootstrap)

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().Bool("json", false,
		assets.FlagDesc(assets.FlagDescKeySystemBootstrapJson),
	)
	cmd.Flags().BoolP("quiet", "q", false,
		assets.FlagDesc(assets.FlagDescKeySystemBootstrapQuiet),
	)

	return cmd
}
