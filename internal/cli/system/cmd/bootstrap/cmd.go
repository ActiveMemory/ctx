//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx system bootstrap" subcommand.
//
// Returns:
//   - *cobra.Command: Configured bootstrap subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemBootstrap)

	c := &cobra.Command{
		Use:   cmd.UseSystemBootstrap,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	c.Flags().Bool(cflag.JSON, false,
		desc.Flag(flag.DescKeySystemBootstrapJson),
	)
	c.Flags().BoolP(cflag.Quiet, cflag.ShortQuiet, false,
		desc.Flag(flag.DescKeySystemBootstrapQuiet),
	)

	return c
}
