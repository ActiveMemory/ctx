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

// deploySkills creates .github/skills/<name>/SKILL.md for each
// embedded Copilot CLI skill template. Skips skills that already exist.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if skill reading or file write fails
func deploySkills(cmd *cobra.Command) error {
	skills, readErr := agent.CopilotCLISkills()
	if readErr != nil {
		return readErr
	}

	skillsBase := filepath.Join(cfgHook.DirGitHub, cfgHook.DirGitHubSkills)
	for name, content := range skills {
		skillDir := filepath.Join(skillsBase, name)
		target := filepath.Join(skillDir, cfgHook.FileSKILLMd)

		if _, statErr := os.Stat(target); statErr == nil {
			writeSetup.InfoCopilotCLISkipped(cmd, target)
			continue
		}

		if mkErr := ctxIo.SafeMkdirAll(skillDir, fs.PermExec); mkErr != nil {
			return mkErr
		}
		if wErr := ctxIo.SafeWriteFile(target, content, fs.PermFile); wErr != nil {
			return wErr
		}
		writeSetup.InfoCopilotCLICreated(cmd, target)
	}
	return nil
}
