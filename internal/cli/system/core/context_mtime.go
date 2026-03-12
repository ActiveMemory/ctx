//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
)

// GetLatestContextMtime returns the most recent mtime of any .context/*.md file.
//
// Parameters:
//   - contextDir: Path to the context directory
//
// Returns:
//   - int64: Unix timestamp of the most recent modification, or 0
func GetLatestContextMtime(contextDir string) int64 {
	entries, readErr := os.ReadDir(contextDir)
	if readErr != nil {
		return 0
	}

	var latest int64
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), file.ExtMarkdown) {
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
