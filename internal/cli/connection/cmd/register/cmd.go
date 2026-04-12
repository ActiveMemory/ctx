//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package register

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreReg "github.com/ActiveMemory/ctx/internal/cli/connection/core/register"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the connect register subcommand.
//
// Returns:
//   - *cobra.Command: The register subcommand
func Cmd() *cobra.Command {
	var adminToken string

	short, long := desc.Command(cmd.DescKeyConnectionRegister)

	c := &cobra.Command{
		Use:     cmd.UseConnectionRegister,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyConnectionRegister),
		Args:    cobra.ExactArgs(1),
		RunE: func(
			cobraCmd *cobra.Command, args []string,
		) error {
			return coreReg.Run(
				cobraCmd, args[0], adminToken,
			)
		},
	}

	flagbind.StringFlag(
		c, &adminToken,
		cFlag.Token, flag.DescKeyConnectionToken,
	)
	_ = c.MarkFlagRequired(cFlag.Token)

	return c
}
