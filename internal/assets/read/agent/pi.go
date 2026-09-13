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

// pi.go is a deliberate deviation from agent.go's layout:
// the Pi integration embeds two surfaces (an extension
// directory and a skill tree), so both accessors live in
// this file instead of joining the single-payload
// accessors in agent.go.

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

// PiSkillReferences reads the embedded reference files of every Pi
// skill. Keys are skill names; values map a reference file name to
// its content. Skills without a references directory are absent
// from the map.
//
// Returns:
//   - map[string]map[string][]byte: Skill -> reference file -> content
//   - error: Non-nil if a read fails
func PiSkillReferences() (map[string]map[string][]byte, error) {
	refs := make(map[string]map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, asset.DirIntegrationsPiSkill)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		refDir := path.Join(
			asset.DirIntegrationsPiSkill, name, asset.DirReferences)
		refEntries, refErr := fs.ReadDir(assets.FS, refDir)
		if refErr != nil {
			// No references directory for this skill.
			continue
		}
		for _, ref := range refEntries {
			if ref.IsDir() {
				continue
			}
			content, readErr := assets.FS.ReadFile(
				path.Join(refDir, ref.Name()))
			if readErr != nil {
				return nil, readErr
			}
			if refs[name] == nil {
				refs[name] = make(map[string][]byte)
			}
			refs[name][ref.Name()] = content
		}
	}
	return refs, nil
}
