//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cursor generates Cursor MCP and steering configuration files.
package cursor

import (
	"github.com/spf13/cobra"

	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Cursor configuration paths.
const (
	// dirCursor is the Cursor editor config directory.
	dirCursor = ".cursor"
	// fileMCPJSON is the MCP server config file name.
	fileMCPJSON = "mcp.json"
	// displayName is the display name for Cursor.
	displayName = "Cursor"
	// mcpConfigPath is the deployed MCP config path.
	mcpConfigPath = dirCursor + "/mcp.json"
	// steeringPath is the deployed steering path.
	steeringPath = dirCursor + "/rules/"
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
	writeSetup.DeployComplete(cmd, displayName, mcpConfigPath, steeringPath)
	return nil
}
