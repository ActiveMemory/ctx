//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cline generates Cline MCP and steering configuration files.
package cline

import (
	"github.com/spf13/cobra"

	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Cline deploy constants.
const (
	// displayName is the display name for Cline.
	displayName = "Cline"
	// mcpConfigPath is the deployed MCP config path.
	mcpConfigPath = ".vscode/mcp.json"
	// steeringPath is the deployed steering path.
	steeringPath = ".clinerules/"
)

// Deploy generates Cline integration files:
//  1. .vscode/mcp.json — MCP server configuration (shared with VS Code)
//  2. .clinerules/*.md — synced steering files
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
