//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the mcp command group.
//
// Returns:
//   - *cobra.Command: The mcp command with subcommands registered
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyMcp, cmd.UseMcp,
		serveCmd(),
	)
}
