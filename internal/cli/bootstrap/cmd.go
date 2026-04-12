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
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx bootstrap" top-level command.
//
// Returns:
//   - *cobra.Command: Configured bootstrap command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyBootstrap)

	c := &cobra.Command{
		Use:     cmd.UseBootstrap,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyBootstrap),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}

	flagbind.BoolFlagNoPtr(c,
		cFlag.JSON, flag.DescKeyBootstrapJson,
	)
	flagbind.BoolFlagShort(c,
		cFlag.Quiet, cFlag.ShortQuiet,
		flag.DescKeyBootstrapQuiet,
	)

	return c
}
