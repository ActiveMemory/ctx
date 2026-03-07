//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/config/cmd/schema"
	"github.com/ActiveMemory/ctx/internal/cli/config/cmd/status"
	"github.com/ActiveMemory/ctx/internal/cli/config/cmd/switchcmd"
)

// Cmd returns the "ctx config" parent command.
//
// Returns:
//   - *cobra.Command: Configured config command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage runtime configuration",
		Long: `Manage runtime configuration profiles.

Subcommands:
  switch [dev|base]    Switch .ctxrc profile (no arg = toggle)
  status               Show active .ctxrc profile
  schema               Print JSON Schema for .ctxrc`,
	}

	cmd.AddCommand(
		switchcmd.Cmd(),
		status.Cmd(),
		schema.Cmd(),
	)

	return cmd
}
