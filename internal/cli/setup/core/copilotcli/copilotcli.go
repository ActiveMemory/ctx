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
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeErr "github.com/ActiveMemory/ctx/internal/write/err"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Deploy generates .github/hooks/ctx-hooks.json for GitHub
// Copilot CLI integration.
//
// The manifest invokes `ctx system ...` commands directly (with a
// repo-root cwd), so no wrapper scripts are shipped. Also writes
// .github/agents/ctx.md and
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
	targetJSON := filepath.Join(hooksDir, cfgHook.FileCopilotCLIHooksJSON)

	// Check if ctx-hooks.json already exists
	if _, statErr := os.Stat(targetJSON); statErr == nil {
		writeSetup.InfoCopilotCLISkipped(cmd, targetJSON)
		return nil
	}

	// Create directories
	if mkErr := ctxIo.SafeMkdirAll(hooksDir, fs.PermExec); mkErr != nil {
		return errFs.Mkdir(hooksDir, mkErr)
	}

	// Write ctx-hooks.json
	jsonContent, readErr := agent.CopilotCLIHooksJSON()
	if readErr != nil {
		return readErr
	}
	wErr := ctxIo.SafeWriteFile(targetJSON, jsonContent, fs.PermFile)
	if wErr != nil {
		return errFs.FileWrite(targetJSON, wErr)
	}
	writeSetup.InfoCopilotCLICreated(cmd, targetJSON)

	// Write .github/agents/ctx.md
	if agentErr := deployGithubAsset(
		cmd,
		cfgHook.DirGitHubAgents, cfgHook.FileAgentsCtxMd,
		agent.AgentsCtxMd,
	); agentErr != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubAgents, agentErr)
	}

	// Write .github/instructions/context.instructions.md
	if instrErr := deployGithubAsset(
		cmd,
		cfgHook.DirGitHubInstructions, cfgHook.FileInstructionsCtxMd,
		agent.InstructionsCtxMd,
	); instrErr != nil {
		writeErr.WarnFile(
			cmd, cfgHook.DirGitHubInstructions, instrErr,
		)
	}

	// Register ctx MCP server in ~/.copilot/mcp-config.json
	if mcpErr := ensureMCPConfig(cmd); mcpErr != nil {
		writeErr.WarnFile(
			cmd, cfgHook.FileMCPConfigJSON, mcpErr,
		)
	}

	// Write .github/skills/<name>/SKILL.md for Copilot CLI skills
	if skillErr := deploySkills(cmd); skillErr != nil {
		writeErr.WarnFile(cmd, cfgHook.DirGitHubSkills, skillErr)
	}

	writeSetup.InfoCopilotCLISummary(cmd)
	return nil
}
