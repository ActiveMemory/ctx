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
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deployAgent creates .github/agents/ctx.md for Copilot CLI custom
// agent delegation. Skips if the file already exists.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func deployAgent(cmd *cobra.Command) error {
	agentsDir := filepath.Join(cfgHook.DirGitHub, cfgHook.DirGitHubAgents)
	target := filepath.Join(agentsDir, cfgHook.FileAgentsCtxMd)

	if _, statErr := os.Stat(target); statErr == nil {
		writeSetup.InfoCopilotCLISkipped(cmd, target)
		return nil
	}

	if mkErr := ctxIo.SafeMkdirAll(agentsDir, fs.PermExec); mkErr != nil {
		return mkErr
	}

	content, readErr := agent.AgentsCtxMd()
	if readErr != nil {
		return readErr
	}
	if wErr := ctxIo.SafeWriteFile(target, content, fs.PermFile); wErr != nil {
		return wErr
	}
	writeSetup.InfoCopilotCLICreated(cmd, target)
	return nil
}
