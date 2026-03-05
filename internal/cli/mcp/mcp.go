//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp provides the CLI command for running the MCP server.
package mcp

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	internalmcp "github.com/ActiveMemory/ctx/internal/mcp"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the mcp command group.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Model Context Protocol server",
		Long:  "Run ctx as an MCP server over stdin/stdout.\n\nThe MCP server exposes context files as resources and ctx commands as tools,\nenabling any MCP-compatible AI tool to access project context.",
	}

	cmd.AddCommand(serveCmd())

	return cmd
}

// serveCmd returns the mcp serve subcommand.
func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "serve",
		Short:        "Start the MCP server (stdin/stdout)",
		Long:         "Start the MCP server, communicating via JSON-RPC 2.0 over stdin/stdout.\n\nThis command is intended to be invoked by MCP clients (AI tools), not\nrun directly by users. Configure your AI tool to run 'ctx mcp serve'\nas an MCP server.",
		Annotations:  map[string]string{config.AnnotationSkipInit: "true"},
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			srv := internalmcp.NewServer(rc.ContextDir())
			return srv.Serve()
		},
	}
}
