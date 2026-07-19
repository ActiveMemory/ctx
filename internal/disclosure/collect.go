//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	"os"
	"path/filepath"
	"strings"

	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
	cfgFile "github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/heading"
)

// gistBasenames returns the set of theme-file basenames the root's gists
// link to.
//
// Parameters:
//   - root: the parsed root
//
// Returns:
//   - map[string]bool: set of basenames (e.g. "hooks.md")
func gistBasenames(root Root) map[string]bool {
	names := map[string]bool{}
	for _, t := range root.Themes {
		if t.Link != "" {
			names[filepath.Base(t.Link)] = true
		}
	}
	return names
}

// themeFiles returns the full paths of the .md files in themeDir. A
// missing directory is not an error — it means the root is not yet
// migrated, so there are no theme files.
//
// Parameters:
//   - themeDir: directory to scan
//
// Returns:
//   - []string: full paths of theme files (nil when the dir is absent)
//   - error: a read error other than "does not exist"
func themeFiles(themeDir string) ([]string, error) {
	entries, err := os.ReadDir(themeDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var paths []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), cfgFile.ExtMarkdown) {
			paths = append(paths, filepath.Join(themeDir, e.Name()))
		}
	}
	return paths, nil
}

// entryIDs returns the identity of every "## [ts] Title" entry in
// content, in order. Identity is timestamp+title, not timestamp alone:
// two entries added in the same second share a timestamp but are
// distinct entries (observed in LEARNINGS.md).
//
// Parameters:
//   - content: markdown to scan for entry blocks
//
// Returns:
//   - []string: entry identities, nil if none
func entryIDs(content string) []string {
	blocks := heading.ParseEntryBlocks(content)
	if len(blocks) == 0 {
		return nil
	}
	ids := make([]string, 0, len(blocks))
	for _, b := range blocks {
		ids = append(ids, b.Entry.Timestamp+cfgDisc.IDSeparator+b.Entry.Title)
	}
	return ids
}

// entryID is a single staged entry's identity: timestamp joined to title
// by IDSeparator, matching the ids in entryIDs and the keys SplitStaging
// returns. Identity is timestamp+title, not timestamp alone: two entries
// added in the same second share a timestamp but are distinct.
//
// Parameters:
//   - e: a staged entry
//
// Returns:
//   - string: the entry's identity
func entryID(e StagedEntry) string {
	return e.Timestamp + cfgDisc.IDSeparator + e.Title
}
