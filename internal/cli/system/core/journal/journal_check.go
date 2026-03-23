//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// NewestMtime returns the most recent mtime (as Unix timestamp) of files
// with the given extension in the directory. Returns 0 if none is found.
//
// Parameters:
//   - dir: absolute path to the directory to scan
//   - ext: file extension to match (e.g., file.ExtMarkdown)
//
// Returns:
//   - int64: Unix timestamp of the newest matching file, or 0
func NewestMtime(dir, ext string) int64 {
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return 0
	}

	var latest int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ext) {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		mtime := info.ModTime().Unix()
		if mtime > latest {
			latest = mtime
		}
	}
	return latest
}

// CountNewerFiles recursively counts files with the given extension that
// are newer than the reference timestamp.
//
// Parameters:
//   - dir: absolute path to the root directory to walk
//   - ext: file extension to match (e.g. ".jsonl")
//   - refTime: Unix timestamp threshold; only files newer than this are counted
//
// Returns:
//   - int: number of matching files newer than refTime
func CountNewerFiles(dir, ext string, refTime int64) int {
	count := 0
	_ = filepath.Walk(dir, func(_ string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil // skip errors
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ext) {
			return nil
		}
		if info.ModTime().Unix() > refTime {
			count++
		}
		return nil
	})
	return count
}

// CountUnenriched counts journal .md files that lack an enriched date
// in the journal state file.
//
// Parameters:
//   - dir: absolute path to the journal directory
//
// Returns:
//   - int: number of unenriched journal entries
func CountUnenriched(dir string) int {
	jstate, loadErr := state.Load(dir)
	if loadErr != nil {
		return 0
	}
	return jstate.CountUnenriched(dir)
}
