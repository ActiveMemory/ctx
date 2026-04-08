//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package shared

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// sharedDir is the subdirectory for shared entries.
const sharedDir = "shared"

// LoadBodies reads all markdown files from .context/shared/
// and returns their contents as strings.
//
// Returns nil if the shared directory does not exist or is
// empty (shared knowledge is opt-in).
//
// Returns:
//   - []string: file contents, one per shared file
func LoadBodies() []string {
	dir := filepath.Join(rc.ContextDir(), sharedDir)
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil
	}

	var bodies []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), file.ExtMarkdown) {
			continue
		}
		data, loadErr := io.SafeReadUserFile(
			filepath.Join(dir, e.Name()),
		)
		if loadErr != nil || len(data) == 0 {
			continue
		}
		bodies = append(bodies, string(data))
	}
	return bodies
}
