//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cursor generates Cursor MCP and steering configuration files.
package cursor

import (
	"github.com/spf13/cobra"

	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates Cursor integration files:
//  1. .cursor/mcp.json — MCP server configuration
//  2. .cursor/rules/*.mdc — synced steering files
func Deploy(cmd *cobra.Command) error {
	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		return mcpErr
	}
	if steerErr := syncSteering(cmd); steerErr != nil {
		return steerErr
	}
	writeSetup.DeployComplete(
		cmd, cfgSetup.DisplayCursor,
		cfgSetup.MCPConfigPathCursor,
		cfgSetup.SteeringPathCursor,
	)
	return nil
}
