//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// SkillList returns available skill directory names.
//
// Each skill is a directory containing a SKILL.md file following the
// Agent Skills specification (https://agentskills.io/specification).
//
// Returns:
//   - []string: List of skill directory names in claude/skills/
//   - error: Non-nil if directory read fails
func SkillList() ([]string, error) {
	entries, err := assets.FS.ReadDir(asset.DirClaudeSkills)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// SkillContent reads a skill's SKILL.md file by skill name.
//
// Parameters:
//   - name: Skill directory name (e.g., "ctx-status")
//
// Returns:
//   - []byte: SKILL.md content from claude/skills/<name>/
//   - error: Non-nil if the file not found or read fails
func SkillContent(name string) ([]byte, error) {
	return assets.FS.ReadFile(path.Join(asset.DirClaudeSkills, name, asset.FileSKILLMd))
}

// SkillReference reads a reference file from a skill's references/ directory.
//
// Parameters:
//   - skill: Skill directory name (e.g., "ctx-skill-audit")
//   - filename: Reference filename (e.g., "anthropic-best-practices.md")
//
// Returns:
//   - []byte: Reference file content
//   - error: Non-nil if the file is not found or read fails
func SkillReference(skill, filename string) ([]byte, error) {
	return assets.FS.ReadFile(path.Join(
		asset.DirClaudeSkills, skill, asset.DirReferences, filename,
	))
}

// SkillReferenceList returns available reference filenames for a skill.
//
// Parameters:
//   - skill: Skill directory name (e.g., "ctx-skill-audit")
//
// Returns:
//   - []string: List of reference filenames
//   - error: Non-nil if the references directory is not found or read fails
func SkillReferenceList(skill string) ([]string, error) {
	entries, readErr := assets.FS.ReadDir(path.Join(
		asset.DirClaudeSkills, skill, asset.DirReferences,
	))
	if readErr != nil {
		return nil, readErr
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
