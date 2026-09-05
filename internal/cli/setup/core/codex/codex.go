//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/codex"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates all project-local Codex integration files.
//
// Writes .codex/hooks.json (create or merge), .codex/config.toml
// (create or append the [mcp_servers.ctx] table), AGENTS.md
// (marker merge), and .agents/skills/<name>/SKILL.md for every
// embedded Codex skill. When the ctx Codex plugin is enabled in
// the user's config.toml, only AGENTS.md is deployed.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if the hooks manifest cannot be written
//     (other errors are warned but do not halt deployment)
func Deploy(cmd *cobra.Command) error {
	home := codex.Home()
	if codex.PluginEnabled(home) {
		if codex.PluginNativeVariant(home) {
			writeSetup.InfoCodexPluginActive(cmd)
			if codex.ProjectConfigured() {
				// A previously deployed project route would
				// make every hook run twice next to the
				// plugin; tell the user which file to remove.
				writeSetup.InfoCodexProjectAlso(cmd)
			}
			deployAgents(cmd)
			writeSetup.InfoCodexSummaryPlugin(cmd)
			return nil
		}
		// A plugin is enabled but the cached copy is the legacy
		// Claude Code variant (installed from a revision without
		// the Codex marketplace). Its hooks cannot run under
		// Codex, so deploy the project-local route in full.
		writeSetup.InfoCodexPluginWrongVariant(cmd)
	}

	hooksOK, hooksErr := deployHooks(cmd)
	if hooksErr != nil {
		return hooksErr
	}

	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		writeErr.WarnFile(cmd, cfgSetup.MCPConfigPathCodex, mcpErr)
	}

	deployAgents(cmd)

	if skillErr := deploySkills(cmd); skillErr != nil {
		writeErr.WarnFile(cmd, cfgSetup.SkillsPathCodex, skillErr)
	}

	if hooksOK {
		writeSetup.InfoCodexSummary(cmd)
	}
	return nil
}
