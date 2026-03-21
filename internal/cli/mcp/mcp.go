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
	"github.com/ActiveMemory/ctx/internal/cli/mcp/cmd/root"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the mcp command group.
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyMcp)
	c := &cobra.Command{
		Use:   cmd.UseMcp,
		Short: short,
		Long:  long,
	}

	c.AddCommand(serveCmd())

	return c
}

// serveCmd returns the mcp serve subcommand.
func serveCmd() *cobra.Command {
	serveShort, serveLong := desc.CommandDesc(cmd.DescKeyMcpServe)
	return &cobra.Command{
		Use:          cmd.UseMcpServe,
		Short:        serveShort,
		Long:         serveLong,
		Annotations:  map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		SilenceUsage: true,
		RunE:         root.Cmd,
	}
}
