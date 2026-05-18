//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package closeout

import (
	"path/filepath"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	errCloseout "github.com/ActiveMemory/ctx/internal/err/closeout"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Archive moves closeout files from their current location into
// the supplied archive directory. The move is performed via
// rename when source and destination share a filesystem; the
// caller is responsible for ensuring archiveDir is writable.
//
// Files are append-never-rewrite; archival physically relocates
// the bytes without modifying them. Archived closeouts are
// immutable.
//
// Parameters:
//   - archiveDir: absolute path to .context/archive/closeouts/
//     (created if absent).
//   - files: closeouts to move (paths must be writable).
//
// Returns:
//   - error: non-nil on mkdir / rename failure. Partial-
//     success is possible: files moved before the failure
//     remain at their archive paths; the error names the file
//     that failed.
func Archive(archiveDir string, files []File) error {
	if len(files) == 0 {
		return nil
	}
	if mkErr := io.SafeMkdirAll(archiveDir, cfgFs.PermExec); mkErr != nil {
		return errCloseout.MkdirArchive(mkErr)
	}
	for _, f := range files {
		dst := filepath.Join(archiveDir, filepath.Base(f.Path))
		if renameErr := io.SafeRename(f.Path, dst); renameErr != nil {
			return errCloseout.ArchiveMove(filepath.Base(f.Path), renameErr)
		}
	}
	return nil
}
