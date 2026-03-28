//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package lookup

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

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

// loadPermissions reads an embedded permission file and splits it into entries.
//
// Parameters:
//   - path: Embedded filesystem path to the permission file
//
// Returns:
//   - []string: Non-empty, non-comment lines from the file; nil on read failure
func loadPermissions(path string) []string {
	data, readErr := assets.FS.ReadFile(path)
	if readErr != nil {
		return nil
	}
	var result []string
	for _, line := range strings.Split(string(data), token.NewlineLF) {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, token.PrefixHeading) {
			continue
		}
		result = append(result, line)
	}
	return result
}
