//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package kiro generates Kiro MCP and steering configuration files.
package kiro

import (
	"github.com/spf13/cobra"

	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
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
		cmd, cfgSetup.DisplayKiro,
		cfgSetup.MCPConfigPathKiro,
		cfgSetup.SteeringDeployPathKiro,
	)
	return nil
}
