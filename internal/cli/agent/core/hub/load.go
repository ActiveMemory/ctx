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
)

// LoadBodies reads all markdown files from .context/hub/
// and returns their contents as strings.
//
// ctxDir is supplied by the caller so this function does not
// re-resolve it; the caller decides whether "no context dir" is
// benign and handles it before invoking us.
//
// Any directory read failure (including a missing hub directory)
// is propagated so the caller can surface it. [LoadBodies] is only
// invoked when the user explicitly requested shared content (e.g.
// `ctx agent --include-share`); telling them "everything is fine,
// here's an empty list" when the hub directory does not exist hides
// a real setup gap.
//
// Per-file read failures inside an existing hub directory are still
// tolerated silently. One unreadable sibling should not blank the
// rest.
//
// Parameters:
//   - ctxDir: absolute path to the context directory
//
// Returns:
//   - []string: file contents, one per shared file
//   - error: non-nil on any directory read failure
func LoadBodies(ctxDir string) ([]string, error) {
	dir := filepath.Join(ctxDir, cfgHub.DirHub)
	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return nil, readErr
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
	return bodies, nil
}
