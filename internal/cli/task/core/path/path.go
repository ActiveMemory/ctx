//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package path

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// FilePath returns the path to TASKS.md.
//
// Returns:
//   - string: Full path to .context/TASKS.md
func FilePath() string {
	return filepath.Join(rc.ContextDir(), ctx.Task)
}

// ArchiveDir returns the path to the archive directory.
//
// Returns:
//   - string: Full path to .context/archive/
func ArchiveDir() string {
	return filepath.Join(rc.ContextDir(), dir.Archive)
}
