//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"errors"
	"os"
	"path/filepath"

	cfgSkill "github.com/ActiveMemory/ctx/internal/config/skill"
	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// LoadAll reads all installed skills from subdirectories of skillsDir.
// Each subdirectory must contain a SKILL.md file with YAML frontmatter.
// Returns an empty slice without error if the skills directory does not exist.
//
// Parameters:
//   - skillsDir: root directory containing skill subdirectories.
//
// Returns:
//   - []*Skill: loaded skills, or nil if the directory is absent.
//   - error: non-nil if a directory entry cannot be read or parsed.
func LoadAll(skillsDir string) ([]*Skill, error) {
	entries, readErr := os.ReadDir(skillsDir)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errSkill.ReadDir(skillsDir, readErr)
	}

	var skills []*Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		sk, loadErr := Load(skillsDir, entry.Name())
		if loadErr != nil {
			return nil, loadErr
		}
		skills = append(skills, sk)
	}
	return skills, nil
}

// Load reads a single skill by name from the given skills directory.
// The name corresponds to a subdirectory containing a SKILL.md file.
//
// Parameters:
//   - skillsDir: root directory containing skill subdirectories.
//   - name: name of the skill subdirectory to load.
//
// Returns:
//   - *Skill: the parsed skill metadata.
//   - error: non-nil if the manifest cannot be read or parsed.
func Load(skillsDir, name string) (*Skill, error) {
	dir := filepath.Join(skillsDir, name)
	manifestPath := filepath.Join(dir, cfgSkill.SkillManifest)

	data, readErr := ctxIo.SafeReadUserFile(manifestPath)
	if readErr != nil {
		return nil, errSkill.Load(name, readErr)
	}

	sk, parseErr := parseManifest(data, name, dir)
	if parseErr != nil {
		return nil, parseErr
	}
	return sk, nil
}
