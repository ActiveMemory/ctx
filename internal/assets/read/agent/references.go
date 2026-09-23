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

// SkillReferences reads the embedded reference files of every skill
// in one skill tree. Keys are skill names; values map a reference
// file name to its content. Skills without a references directory
// are absent from the map.
//
// One reader serves every tree (canonical and generated): the shape
// is identical, so a per-tool copy only invites the next copy to
// drift.
//
// Parameters:
//   - treeDir: Embedded skill tree root, e.g. asset.DirCodexSkills
//
// Returns:
//   - map[string]map[string][]byte: Skill -> reference file -> content
//   - error: Non-nil if the tree read or a file read fails
func SkillReferences(treeDir string) (map[string]map[string][]byte, error) {
	refs := make(map[string]map[string][]byte)
	entries, dirErr := fs.ReadDir(assets.FS, treeDir)
	if dirErr != nil {
		return nil, dirErr
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		refDir := path.Join(treeDir, name, asset.DirReferences)
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
