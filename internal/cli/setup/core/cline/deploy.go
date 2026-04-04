//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cline

import (
	"github.com/spf13/cobra"

	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	cfgVscode "github.com/ActiveMemory/ctx/internal/config/vscode"

	mcpDeploy "github.com/ActiveMemory/ctx/internal/cli/setup/core/mcp"
)

// ensureMCPConfig creates .vscode/mcp.json with the ctx
// MCP server configuration. Skips if the file exists.
func ensureMCPConfig(cmd *cobra.Command) error {
	cfg := vscodeMCPConfig{
		Servers: map[string]vscodeMCPServer{
			mcpServer.Name: {
				Command: mcpServer.Command,
				Args:    mcpServer.Args(),
			},
		},
	}
	return mcpDeploy.Deploy(
		cmd, cfgVscode.Dir,
		cfgVscode.FileMCPJSON, cfg,
	)
}

// syncSteering syncs steering files to Cline format
// if a steering directory exists.
func syncSteering(cmd *cobra.Command) error {
	return mcpDeploy.SyncSteering(
		cmd, cfgHook.ToolCline,
	)
}
