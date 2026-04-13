//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude_check

import (
	"encoding/json"

	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// shortShaLen is the number of characters of a git commit
// SHA to display in the Ready-state output.
const shortShaLen = 8

// readInstalled loads and parses
// ~/.claude/plugins/installed_plugins.json.
//
// Parameters:
//   - path: absolute path to installed_plugins.json
//
// Returns:
//   - installedPluginsFile: parsed file contents
//   - bool: true iff the file exists and parses cleanly
func readInstalled(path string) (installedPluginsFile, bool) {
	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		return installedPluginsFile{}, false
	}
	var f installedPluginsFile
	if unmarshalErr := json.Unmarshal(data, &f); unmarshalErr != nil {
		return installedPluginsFile{}, false
	}
	return f, true
}

// readMarketplaces loads and parses
// ~/.claude/plugins/known_marketplaces.json.
//
// Parameters:
//   - path: absolute path to known_marketplaces.json
//
// Returns:
//   - knownMarketplacesFile: parsed file contents
//   - bool: true iff the file exists and parses cleanly
func readMarketplaces(path string) (knownMarketplacesFile, bool) {
	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		return nil, false
	}
	var f knownMarketplacesFile
	if unmarshalErr := json.Unmarshal(data, &f); unmarshalErr != nil {
		return nil, false
	}
	return f, true
}

// shortSha returns the first [shortShaLen] characters of a
// git SHA, or the full string if it's shorter.
//
// Parameters:
//   - sha: full git SHA (may be empty)
//
// Returns:
//   - string: shortened SHA or empty
func shortSha(sha string) string {
	if len(sha) <= shortShaLen {
		return sha
	}
	return sha[:shortShaLen]
}
