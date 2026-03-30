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

	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/file"
)

// loadYAML parses an embedded YAML file into a commandEntry map.
//
// Parameters:
//   - p: embedded file path to read
//
// Returns:
//   - map[string]commandEntry: parsed entries, empty map on error
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
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), file.ExtYAML) {
			continue
		}
		for k, v := range loadYAML(path.Join(dir, entry.Name())) {
			merged[k] = v
		}
	}
	return merged
}
