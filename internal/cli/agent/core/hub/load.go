//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// LoadBodies reads all markdown files from .context/hub/
// and returns their contents as strings.
//
// Returns nil if the shared directory does not exist or is
// empty (shared knowledge is opt-in).
//
// Returns:
//   - []string: file contents, one per shared file
func LoadBodies() []string {
	dir := filepath.Join(rc.ContextDir(), cfgHub.DirHub)
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
