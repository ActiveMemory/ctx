//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
)

// ScanCandidates walks upward from start collecting every directory
// whose basename matches the canonical context directory name
// (`.context`). The scan is read-only: it does not resolve, bind, or
// select a context directory. It exists so error messages and the
// `ctx activate` subcommand can share the same candidate enumeration
// without reintroducing walk-up resolution elsewhere.
//
// The scan always uses the canonical `.context` basename, independent
// of any `.ctxrc` configuration. Under the explicit-declaration model,
// a custom name is only ever reached via an explicit --context-dir or
// CTX_DIR, so a rename-aware scan would be surplus machinery.
//
// Parameters:
//   - start: directory to begin the upward walk from; typically the
//     current working directory returned by os.Getwd.
//
// Returns:
//   - []string: absolute paths of every matching directory found,
//     ordered innermost-first (closest to start first). Empty when
//     no candidates are visible on the upward path.
func ScanCandidates(start string) []string {
	var out []string
	cur := start
	for {
		path := filepath.Join(cur, dir.Context)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			out = append(out, path)
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return out
		}
		cur = parent
	}
}
