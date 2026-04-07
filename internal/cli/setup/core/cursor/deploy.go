//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cursor

import (
	"github.com/spf13/cobra"

	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"

	mcpDeploy "github.com/ActiveMemory/ctx/internal/cli/setup/core/mcp"
)

// ensureMCPConfig creates .cursor/mcp.json with the ctx
// MCP server configuration. Skips if the file exists.
//
// Parameters:
//   - cmd: cobra command for status output
//
// Returns:
//   - error: config serialization or write failure
func ensureMCPConfig(cmd *cobra.Command) error {
	cfg := mcpConfig{
		MCPServers: map[string]serverEntry{
			mcpServer.Name: {
				Command: mcpServer.Command,
				Args:    mcpServer.Args(),
			},
		},
	}
	return mcpDeploy.Deploy(
		cmd, cfgSetup.DirCursor,
		cfgSetup.FileMCPJSONCursor, cfg,
	)
}

// syncSteering syncs steering files to Cursor format
// if a steering directory exists.
//
// Parameters:
//   - cmd: cobra command for status output
//
// Returns:
//   - error: steering sync failure
func syncSteering(cmd *cobra.Command) error {
	return mcpDeploy.SyncSteering(
		cmd, cfgHook.ToolCursor,
	)
}
