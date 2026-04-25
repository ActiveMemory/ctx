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

// File returns the path to TASKS.md.
//
// Returns:
//   - string: Full path to .context/TASKS.md
//   - error: non-nil when the context directory is not declared
func File() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, ctx.Task), nil
}

// ArchiveDir returns the path to the archive directory.
//
// Returns:
//   - string: Full path to .context/archive/
//   - error: non-nil when the context directory is not declared
func ArchiveDir() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, dir.Archive), nil
}
