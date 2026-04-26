//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deploySkills creates .opencode/skills/<name>/SKILL.md for each
// embedded OpenCode skill. Skips skills whose SKILL.md already
// exists.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func deploySkills(cmd *cobra.Command) error {
	skills, readErr := agent.OpenCodeSkills()
	if readErr != nil {
		return readErr
	}

	skillsBase := filepath.Join(
		cfgHook.DirOpenCode, cfgHook.DirOpenCodeSkills,
	)

	for name, content := range skills {
		skillDir := filepath.Join(skillsBase, name)
		target := filepath.Join(skillDir, cfgHook.FileSKILLMd)

		if _, statErr := os.Stat(target); statErr == nil {
			writeSetup.InfoOpenCodeSkipped(cmd, target)
			continue
		}

		if mkErr := ctxIo.SafeMkdirAll(
			skillDir, fs.PermExec,
		); mkErr != nil {
			return errFs.Mkdir(skillDir, mkErr)
		}

		if wErr := ctxIo.SafeWriteFile(
			target, content, fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		writeSetup.InfoOpenCodeCreated(cmd, target)
	}

	return nil
}
