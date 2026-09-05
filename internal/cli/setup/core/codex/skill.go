//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/agent"
	"github.com/ActiveMemory/ctx/internal/codex"
	cfgAsset "github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHook "github.com/ActiveMemory/ctx/internal/config/hook"
	cfgSetup "github.com/ActiveMemory/ctx/internal/config/setup"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// deploySkills creates .agents/skills/<name>/SKILL.md for each
// embedded Codex skill. Identical files are skipped, stale
// ctx-managed files (frontmatter `name:` matches the directory)
// are refreshed, and foreign files are rejected with a notice.
//
// Parameters:
//   - cmd: Cobra command for output messages
//
// Returns:
//   - error: Non-nil if directory creation or file write fails
func deploySkills(cmd *cobra.Command) error {
	skills, readErr := agent.CodexSkills()
	if readErr != nil {
		return readErr
	}
	refs, refErr := agent.CodexSkillReferences()
	if refErr != nil {
		return refErr
	}

	// Iterate in sorted order so deploy is deterministic.
	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		content := skills[name]
		skillDir := filepath.Join(cfgSetup.SkillsPathCodex, name)
		target := filepath.Join(skillDir, cfgHook.FileSKILLMd)
		if _, validateErr := validateManagedTarget(target); validateErr != nil {
			return validateErr
		}

		if existing, statErr := ctxIo.SafeReadUserFile(target); statErr == nil {
			if bytes.Equal(existing, content) {
				writeSetup.InfoCodexSkipped(cmd, target)
				if refDeployErr := deployReferences(
					cmd, name, refs[name],
				); refDeployErr != nil {
					return refDeployErr
				}
				continue
			}
			if !codex.SkillManaged(existing, name) {
				writeSetup.InfoCodexRejected(cmd, target)
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
		writeSetup.InfoCodexCreated(cmd, target)
		if refDeployErr := deployReferences(
			cmd, name, refs[name],
		); refDeployErr != nil {
			return refDeployErr
		}
	}

	return nil
}

// deployReferences writes a skill's reference files under its
// deployed directory. Runs only after the skill's SKILL.md was
// deployed or confirmed ctx-managed, so references never land in
// a foreign skill directory. Identical files are left untouched.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - name: skill directory name
//   - files: reference file name -> content (may be nil)
//
// Returns:
//   - error: Non-nil if directory creation or a write fails
func deployReferences(
	cmd *cobra.Command,
	name string,
	files map[string][]byte,
) error {
	if len(files) == 0 {
		return nil
	}
	refDir := filepath.Join(
		cfgSetup.SkillsPathCodex, name, cfgAsset.DirReferences,
	)

	refNames := make([]string, 0, len(files))
	for refName := range files {
		refNames = append(refNames, refName)
	}
	sort.Strings(refNames)

	for _, refName := range refNames {
		target := filepath.Join(refDir, refName)
		if _, validateErr := validateManagedTarget(target); validateErr != nil {
			return validateErr
		}
		content := files[refName]
		if existing, statErr := ctxIo.SafeReadUserFile(target); statErr == nil {
			if bytes.Equal(existing, content) {
				continue
			}
		} else if !os.IsNotExist(statErr) {
			return errFs.FileRead(target, statErr)
		}
		if mkErr := ctxIo.SafeMkdirAll(refDir, fs.PermExec); mkErr != nil {
			return errFs.Mkdir(refDir, mkErr)
		}
		if wErr := ctxIo.SafeWriteFile(
			target, content, fs.PermFile,
		); wErr != nil {
			return errFs.FileWrite(target, wErr)
		}
		writeSetup.InfoCodexCreated(cmd, target)
	}
	return nil
}
