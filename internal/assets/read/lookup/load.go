//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"io/fs"
	"path"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"

	"gopkg.in/yaml.v3"
)

// loadYAML parses an embedded YAML file into a commandEntry map.
func loadYAML(p string) map[string]commandEntry {
	data, readErr := assets.FS.ReadFile(p)
	if readErr != nil {
		return make(map[string]commandEntry)
	}
	m := make(map[string]commandEntry)
	if parseErr := yaml.Unmarshal(data, &m); parseErr != nil {
		return make(map[string]commandEntry)
	}
	return m
}

// loadYAMLDir reads all YAML files in an embedded directory and merges
// them into a single commandEntry map.
//
// Parameters:
//   - dir: embedded directory path to scan
//
// Returns:
//   - map[string]commandEntry: merged entries from all files
func loadYAMLDir(dir string) map[string]commandEntry {
	merged := make(map[string]commandEntry)
	entries, readErr := fs.ReadDir(assets.FS, dir)
	if readErr != nil {
		return merged
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}
		for k, v := range loadYAML(path.Join(dir, entry.Name())) {
			merged[k] = v
		}
	}
	return merged
}

// Init loads all embedded YAML description maps. Call once from main()
// before building the command tree. Tests that need descriptions must
// call Init() in their setup.
func Init() {
	CommandsMap = loadYAML(asset.PathCommandsYAML)
	FlagsMap = loadYAML(asset.PathFlagsYAML)
	TextMap = loadYAMLDir(asset.DirCommandsText)
	ExamplesMap = loadYAML(asset.PathExamplesYAML)
	allowPerms = loadPermissions(asset.PathAllowTxt)
	denyPerms = loadPermissions(asset.PathDenyTxt)
	stopWordsMap = loadStopWords()
}
