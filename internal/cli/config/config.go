//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx config" parent command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage runtime configuration",
		Long: `Manage runtime configuration profiles.

Subcommands:
  switch [dev|base]    Switch .ctxrc profile (no arg = toggle)
  status               Show active .ctxrc profile`,
	}

	cmd.AddCommand(
		switchCmd(),
		statusCmd(),
	)

	return cmd
}
