//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package kiro generates Kiro MCP and steering configuration files.
package kiro

import (
	"github.com/spf13/cobra"

	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Kiro configuration paths.
const (
	// DirKiro is the Kiro editor config directory.
	DirKiro = ".kiro"
	// DirSettings is the Kiro settings subdirectory.
	DirSettings = "settings"
	// FileMCPJSON is the MCP server config file name.
	FileMCPJSON = "mcp.json"
	// displayName is the display name for Kiro.
	displayName = "Kiro"
	// mcpConfigPath is the deployed MCP config path.
	mcpConfigPath = DirKiro + "/settings/mcp.json"
	// steeringDeployPath is the deployed steering path.
	steeringDeployPath = DirKiro + "/steering/"
)

// Deploy generates Kiro integration files:
//  1. .kiro/settings/mcp.json — MCP server configuration
//  2. .kiro/steering/*.md — synced steering files
//
// Skips files that already exist to avoid overwriting user customizations.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func Deploy(cmd *cobra.Command) error {
	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		return mcpErr
	}

	if steerErr := syncSteering(cmd); steerErr != nil {
		return steerErr
	}

	writeSetup.DeployComplete(
		cmd, displayName,
		mcpConfigPath,
		steeringDeployPath,
	)
	return nil
}
