//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// TasksFilePath returns the path to TASKS.md.
//
// Returns:
//   - string: Full path to .context/TASKS.md
func TasksFilePath() string {
	return filepath.Join(rc.ContextDir(), ctx.Task)
}

// ArchiveDirPath returns the path to the archive directory.
//
// Returns:
//   - string: Full path to .context/archive/
func ArchiveDirPath() string {
	return filepath.Join(rc.ContextDir(), dir.Archive)
}
