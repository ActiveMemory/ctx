//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package copilotcli

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates .github/hooks/ctx-hooks.json and the
// accompanying hook scripts for GitHub Copilot CLI integration.
//
// Creates the .github/hooks/ and .github/hooks/scripts/ directories if
// needed and writes the JSON config plus bash and PowerShell scripts
// from embedded assets. Also writes .github/agents/ctx.md and
// .github/instructions/context.instructions.md for Copilot CLI.
// Skips if ctx-hooks.json already exists.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func Deploy(cmd *cobra.Command) error {
	hooksDir := filepath.Join(cfgHook.DirGitHub, cfgHook.DirGitHubHooks)
	scriptsDir := filepath.Join(hooksDir, cfgHook.DirGitHubHooksScripts)
	targetJSON := filepath.Join(hooksDir, cfgHook.FileCopilotCLIHooksJSON)

	// Check if ctx-hooks.json already exists
	if _, err := os.Stat(targetJSON); err == nil {
		writeSetup.InfoCopilotCLISkipped(cmd, targetJSON)
		return nil
	}

	// Create directories
	if err := os.MkdirAll(scriptsDir, fs.PermExec); err != nil {
		return errFs.Mkdir(scriptsDir, err)
	}

	// Write ctx-hooks.json
	jsonContent, readErr := agent.CopilotCLIHooksJSON()
	if readErr != nil {
		return readErr
	}
	if wErr := os.WriteFile(targetJSON, jsonContent, fs.PermFile); wErr != nil {
		return errFs.FileWrite(targetJSON, wErr)
	}
	writeSetup.InfoCopilotCLICreated(cmd, targetJSON)

	// Write all hook scripts
	scripts, scrErr := agent.CopilotCLIScripts()
	if scrErr != nil {
		return scrErr
	}
	for name, content := range scripts {
		target := filepath.Join(scriptsDir, name)
		if wErr := os.WriteFile(target, content, fs.PermExec); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		writeSetup.InfoCopilotCLICreated(cmd, target)
	}

	// Write .github/agents/ctx.md
	if err := deployAgent(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubAgents, err)
	}

	// Write .github/instructions/context.instructions.md
	if err := deployInstructions(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubInstructions, err)
	}

	// Register ctx MCP server in ~/.copilot/mcp-config.json
	if err := ensureMCPConfig(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgHook.FileMCPConfigJSON, err)
	}

	// Write .github/skills/<name>/SKILL.md for Copilot CLI skills
	if err := deploySkills(cmd); err != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubSkills, err)
	}

	writeSetup.InfoCopilotCLISummary(cmd)
	return nil
}
