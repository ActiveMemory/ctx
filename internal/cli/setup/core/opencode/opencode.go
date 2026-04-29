//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"github.com/spf13/cobra"

	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates all OpenCode integration files.
//
// Creates the .opencode/plugins/ctx.ts plugin file, registers
// the ctx MCP server in opencode.json, deploys AGENTS.md with
// shared instructions, and copies ctx skills to
// .opencode/skills/.
//
// Skips existing files (idempotent).
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if plugin deployment fails (other errors are
//     warned but do not halt deployment)
func Deploy(cmd *cobra.Command) error {
	if pluginErr := deployPlugin(cmd); pluginErr != nil {
		return pluginErr
	}

	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		writeErr.WarnFile(
			cmd, cfgSetup.MCPConfigPathOpenCode, mcpErr,
		)
	}

	if agentsErr := deployAgents(cmd); agentsErr != nil {
		writeErr.WarnFile(
			cmd, cfgHook.FileAgentsMd, agentsErr,
		)
	}

	if skillErr := deploySkills(cmd); skillErr != nil {
		writeErr.WarnFile(
			cmd, cfgSetup.SkillsPathOpenCode, skillErr,
		)
	}

	writeSetup.InfoOpenCodeSummary(cmd)
	return nil
}
