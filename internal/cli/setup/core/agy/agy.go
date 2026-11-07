//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/agy"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates all project-local Antigravity integration files.
//
// Writes .agents/hooks.json (create or merge named group "ctx"),
// ~/.gemini/config/mcp_config.json (merge-safe MCP server registration),
// AGENTS.md (marker merge), and .agents/skills/<name>/SKILL.md (+ references/)
// for every embedded Antigravity skill. When the ctx Antigravity plugin is
// enabled in ~/.gemini/config/plugins/ctx/, only AGENTS.md is deployed.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if the hooks manifest cannot be written
//     (other errors are warned but do not halt deployment)
func Deploy(cmd *cobra.Command) error {
	configDir := agy.ConfigDir()
	if agy.PluginEnabled(configDir) {
		writeSetup.InfoAgyPluginActive(cmd)
		if agy.ProjectConfigured() {
			writeSetup.InfoAgyProjectAlso(cmd)
		}
		deployAgents(cmd)
		writeSetup.InfoAgySummaryPlugin(cmd)
		return nil
	}

	hooksOK, hooksErr := deployHooks(cmd)
	if hooksErr != nil {
		return hooksErr
	}

	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		target := filepath.Join(configDir, cfgAsset.FileMCPConfigJSON)
		writeErr.WarnFile(cmd, target, mcpErr)
	}

	deployAgents(cmd)

	if skillErr := deploySkills(cmd); skillErr != nil {
		writeErr.WarnFile(cmd, cfgSetup.SkillsPathAgy, skillErr)
	}

	if hooksOK {
		writeSetup.InfoAgySummary(cmd)
	}
	return nil
}
