//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx notify setup" subcommand.
//
// Returns:
//   - *cobra.Command: Configured setup subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("notify.setup")
	return &cobra.Command{
		Use:   "setup",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
