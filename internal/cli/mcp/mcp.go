//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp provides the CLI command for running the MCP server.
package mcp

import (
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/mcp/cmd/root"
)

// Cmd returns the mcp command group.
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyMcp)
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(serveCmd())

	return cmd
}

// serveCmd returns the mcp serve subcommand.
func serveCmd() *cobra.Command {
	serveShort, serveLong := assets.CommandDesc(assets.CmdDescKeyMcpServe)
	return &cobra.Command{
		Use:          "serve",
		Short:        serveShort,
		Long:         serveLong,
		Annotations:  map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		SilenceUsage: true,
		RunE:         root.Cmd,
	}
}
