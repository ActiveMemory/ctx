//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/memory"
)

// statePath returns the filesystem path to the memory state JSON file.
//
// Parameters:
//   - contextDir: Root context directory
//
// Returns:
//   - string: Absolute path to the state file
func statePath(contextDir string) string {
	return filepath.Join(contextDir, dir.State, memory.State)
}
