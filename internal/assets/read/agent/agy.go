//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"io/fs"
	"path"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
)

// AgyHooksJSON reads the embedded Antigravity hooks manifest.
//
// Returns:
//   - []byte: JSON content from agy/hooks.json
//   - error: Non-nil if the file is not found or read fails
func AgyHooksJSON() ([]byte, error) {
	return assets.FS.ReadFile(asset.PathAgyHooksJSON)
}

// AgySkills reads all embedded Antigravity skill templates.
// Returns a map of skill directory name to SKILL.md content for
// skills in agy/skills/.
//
// Returns:
//   - map[string][]byte: Skill name -> SKILL.md content
//   - error: Non-nil if the directory read fails
func AgySkills() (map[string][]byte, error) {
	skills := make(map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirAgySkills)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		skillPath := path.Join(
			asset.DirAgySkills, name, asset.FileSKILLMd,
		)
		content, readErr := assets.FS.ReadFile(skillPath)
		if readErr != nil {
			return nil, readErr
		}
		skills[name] = content
	}
	return skills, nil
}
