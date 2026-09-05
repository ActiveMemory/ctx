//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package kiro

import (
	"path/filepath"

	"github.com/spf13/cobra"

	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	mcpServer "github.com/ActiveMemory/ctx/internal/config/mcp/server"
	cfgMcpTool "github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"

	mcpDeploy "github.com/ActiveMemory/ctx/internal/cli/setup/core/mcp"
)

// ensureMCPConfig creates .kiro/settings/mcp.json with
// the ctx MCP server config. Skips if the file exists.
//
// Parameters:
//   - cmd: cobra command for status output
//
// Returns:
//   - error: config serialization or write failure
func ensureMCPConfig(cmd *cobra.Command) error {
	settingsDir := filepath.Join(
		cfgSetup.DirKiro, cfgSetup.DirSettings,
	)
	cfg := mcpConfig{
		MCPServers: map[string]serverEntry{
			mcpServer.Name: {
				// Deliberately the bare binary name: this file is
				// project-scoped (committable, shared across the
				// team), so a machine-specific absolute path must
				// not be embedded. Each machine resolves ctx from
				// its own PATH.
				Command:  mcpServer.Command,
				Args:     mcpServer.Args(),
				Disabled: false,
				AutoApprove: []string{
					cfgMcpTool.Status,
					cfgMcpTool.SteeringGet,
					cfgMcpTool.Search,
					cfgMcpTool.SessionStart,
					cfgMcpTool.SessionEnd,
					cfgMcpTool.Next,
					cfgMcpTool.Remind,
				},
			},
		},
	}
	return mcpDeploy.Deploy(
		cmd, settingsDir,
		cfgSetup.FileMCPJSON, cfg,
	)
}

// syncSteering syncs steering files to Kiro format
// if a steering directory exists.
//
// Parameters:
//   - cmd: cobra command for status output
//
// Returns:
//   - error: steering sync failure
func syncSteering(cmd *cobra.Command) error {
	return mcpDeploy.SyncSteering(
		cmd, cfgHook.ToolKiro,
	)
}
