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
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx config" parent command.
//
// Returns:
//   - *cobra.Command: Configured config command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyConfig, cmd.UseConfig,
		switchcmd.Cmd(),
		status.Cmd(),
		schema.Cmd(),
	)
}
