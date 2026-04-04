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

	errSkill "github.com/ActiveMemory/ctx/internal/err/skill"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// skillManifest is the expected filename inside each skill directory.
const skillManifest = "SKILL.md"

// frontmatterDelimiter is the YAML frontmatter boundary marker.
const frontmatterDelimiter = "---"

// trimCR is the character set trimmed from the start of
// raw frontmatter content to normalize line endings.
const trimCR = "\n\r"

// LoadAll reads all installed skills from subdirectories of skillsDir.
// Each subdirectory must contain a SKILL.md file with YAML frontmatter.
// Returns an empty slice without error if the skills directory does not exist.
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
func Load(skillsDir, name string) (*Skill, error) {
	dir := filepath.Join(skillsDir, name)
	manifestPath := filepath.Join(dir, skillManifest)

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
