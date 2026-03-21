//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system bootstrap" subcommand.
//
// Returns:
//   - *cobra.Command: Configured bootstrap subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemBootstrap)

	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	cmd.Flags().Bool("json", false,
		desc.Flag(flag.DescKeySystemBootstrapJson),
	)
	cmd.Flags().BoolP("quiet", "q", false,
		desc.Flag(flag.DescKeySystemBootstrapQuiet),
	)

	return cmd
}
