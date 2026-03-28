//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp provides the CLI command for running the MCP server.
package mcp

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the mcp command group.
//
// Returns:
//   - *cobra.Command: The mcp command with subcommands registered
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyMcp)
	c := &cobra.Command{
		Use:   cmd.UseMcp,
		Short: short,
		Long:  long,
	}

	c.AddCommand(serveCmd())

	return c
}
