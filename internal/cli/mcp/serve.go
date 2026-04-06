//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/mcp/cmd/root"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// serveCmd returns the mcp serve subcommand.
//
// Returns:
//   - *cobra.Command: Configured serve subcommand with init-skip annotation
func serveCmd() *cobra.Command {
	serveShort, serveLong := desc.Command(cmd.DescKeyMcpServe)
	return &cobra.Command{
		Use:          cmd.UseMcpServe,
		Short:        serveShort,
		Long:         serveLong,
		Example:      desc.Example(cmd.DescKeyMcpServe),
		Annotations:  map[string]string{cli.AnnotationSkipInit: cli.AnnotationTrue},
		SilenceUsage: true,
		RunE:         root.Cmd,
	}
}
