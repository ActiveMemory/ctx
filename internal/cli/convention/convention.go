//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package convention

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/convention/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the convention command with subcommands.
//
// The convention command provides utilities for managing the
// CONVENTIONS.md file, currently limited to adding new
// entries via "ctx convention add".
//
// Returns:
//   - *cobra.Command: The convention command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyConvention, cmd.UseConvention,
		add.Cmd(),
	)
}
