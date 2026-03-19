//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
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
	short, long := desc.CommandDesc(cmd.DescKeyConfig)

	cmd := &cobra.Command{
		Use:   cmd.DescKeyConfig,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(
		switchcmd.Cmd(),
		status.Cmd(),
		schema.Cmd(),
	)

	return cmd
}
