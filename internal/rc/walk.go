//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"os"
	"path/filepath"
)

// walkForContextDir walks upward from the current working directory
// looking for an existing directory whose basename matches name.
//
// Absolute configured names skip the walk entirely. When no matching
// directory is found upward, returns filepath.Join(cwd, name) as an
// absolute path so that ctx init can create a fresh context directory at
// the current location.
//
// Parameters:
//   - name: Configured context directory name (may be relative or absolute)
//
// Returns:
//   - string: Absolute path to the resolved context directory
func walkForContextDir(name string) string {
	if filepath.IsAbs(name) {
		return name
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return name
	}

	cur := cwd
	for {
		candidate := filepath.Join(cur, name)
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}

	return filepath.Join(cwd, name)
}
