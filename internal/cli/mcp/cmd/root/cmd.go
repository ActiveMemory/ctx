//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	internalMcp "github.com/ActiveMemory/ctx/internal/mcp/server"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd starts the MCP server over stdin/stdout.
//
// Parameters:
//   - cmd: Cobra command for version access
//   - _: Unused positional arguments
//
// Returns:
//   - error: Non-nil if the server fails to start or encounters an I/O error
func Cmd(cmd *cobra.Command, _ []string) error {
	ctxDir, err := rc.RequireContextDir()
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}
	srv := internalMcp.New(ctxDir, cmd.Root().Version)
	return srv.Serve()
}
