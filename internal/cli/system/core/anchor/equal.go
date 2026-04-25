//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package anchor

import "path/filepath"

// Equal reports whether a and b refer to the same directory.
// First compares cleaned paths byte-for-byte (the cheap case);
// then falls back to [filepath.EvalSymlinks] resolution to catch
// symlink-equivalent paths like macOS's `/tmp` → `/private/tmp`.
//
// Resolution failure on either side defaults to "different,"
// which is the correct call: an inherited CTX_DIR pointing at a
// deleted directory is genuine drift even if the injected one
// resolves cleanly. The drift hook's job is to surface real
// misalignment; an over-eager symlink fix that swallowed
// resolution failures would silently hide it.
//
// Parameters:
//   - a: first path (typically the parent-shell inherited CTX_DIR).
//   - b: second path (typically the Claude-injected CTX_DIR).
//
// Returns:
//   - bool: true when the two paths resolve to the same directory.
func Equal(a, b string) bool {
	aClean := filepath.Clean(a)
	bClean := filepath.Clean(b)
	if aClean == bClean {
		return true
	}
	aResolved, aErr := filepath.EvalSymlinks(aClean)
	if aErr != nil {
		return false
	}
	bResolved, bErr := filepath.EvalSymlinks(bClean)
	if bErr != nil {
		return false
	}
	return aResolved == bResolved
}
