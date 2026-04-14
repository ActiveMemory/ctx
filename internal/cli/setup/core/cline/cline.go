//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cline

import (
	"github.com/spf13/cobra"

	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates Cline integration files:
//  1. .vscode/mcp.json — MCP server configuration (shared with VS Code)
//  2. .clinerules/*.md — synced steering files
//
// Parameters:
//   - cmd: Cobra command for output and flag access
//
// Returns:
//   - error: Non-nil if MCP config or steering sync fails
func Deploy(cmd *cobra.Command) error {
	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		return mcpErr
	}
	if steerErr := syncSteering(cmd); steerErr != nil {
		return steerErr
	}
	writeSetup.DeployComplete(
		cmd, cfgSetup.DisplayCline,
		cfgSetup.MCPConfigPathCline,
		cfgSetup.SteeringPathCline,
	)
	return nil
}
