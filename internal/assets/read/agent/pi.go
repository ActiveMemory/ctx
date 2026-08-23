//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"io/fs"
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// PiExtension reads all embedded Pi extension files.
// Returns a map of filename to content for files in
// integrations/pi/extension/.
//
// Returns:
//   - map[string][]byte: Filename -> content for each extension file
//   - error: Non-nil if the directory read fails
func PiExtension() (map[string][]byte, error) {
	files := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(
		assets.FS, asset.DirIntegrationsPiExtension)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		p := path.Join(asset.DirIntegrationsPiExtension, name)
		content, readErr := assets.FS.ReadFile(p)
		if readErr != nil {
			return nil, readErr
		}
		files[name] = content
	}
	return files, nil
}

// PiSkills reads all embedded Pi skill templates.
// Returns a map of skill directory name to SKILL.md content for skills
// in integrations/pi/skills/.
//
// Returns:
//   - map[string][]byte: Skill name -> SKILL.md content
//   - error: Non-nil if the directory read fails
func PiSkills() (map[string][]byte, error) {
	skills := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirIntegrationsPiSkill)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillPath := path.Join(
			asset.DirIntegrationsPiSkill,
			name, asset.FileSKILLMd)
		content, readErr := assets.FS.ReadFile(skillPath)
		if readErr != nil {
			return nil, readErr
		}
		skills[name] = content
	}
	return skills, nil
}
