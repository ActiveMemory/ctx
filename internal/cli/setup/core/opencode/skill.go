//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package opencode

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deploySkills creates .opencode/skills/<name>/SKILL.md for each
// embedded OpenCode skill. Existing ctx-managed skill files are
// refreshed when their contents differ from the embedded source.
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

	// Iterate in sorted order so deploy is deterministic and tests
	// that plant blocking files at a specific skill path observe a
	// stable filesystem state on partial-failure paths.
	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		content := skills[name]
		skillDir := filepath.Join(skillsBase, name)
		target := filepath.Join(skillDir, cfgHook.FileSKILLMd)
		if _, validateErr := validateManagedTarget(target); validateErr != nil {
			return validateErr
		}

		if existing, statErr := ctxIo.SafeReadUserFile(target); statErr == nil {
			if bytes.Equal(existing, content) {
				writeSetup.InfoOpenCodeSkipped(cmd, target)
				continue
			}
		} else if !os.IsNotExist(statErr) {
			return errFs.FileRead(target, statErr)
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
