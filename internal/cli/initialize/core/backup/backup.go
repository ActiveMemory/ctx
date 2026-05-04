//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup writes a timestamped snapshot of the populated
// essential context files before ctx init --reset overwrites
// them.
//
// The backup directory lives inside the same .context/ tree
// (.context/.backup-init-<UTC-ISO>/) so it travels with the
// project and is easy to recover from. The leading dot keeps
// it out of glob-driven listings; the timestamp suffix makes
// every reset uniquely identifiable.
package backup

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgInit "github.com/ActiveMemory/ctx/internal/config/initialize"
	errInit "github.com/ActiveMemory/ctx/internal/err/initialize"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// WriteSnapshot copies each named file from contextDir into a
// fresh timestamped backup directory and returns the absolute
// path of that directory.
//
// Best-effort: a file that has disappeared between the
// populated-files probe and the snapshot call is skipped (no
// error). Any other read or write failure aborts and surfaces
// the cause to the caller; init must not proceed with an
// incomplete backup.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//   - files:      basenames (relative to contextDir) to copy
//
// Returns:
//   - string: absolute path of the backup directory created
//   - error:  non-nil on mkdir, read, or write failure
func WriteSnapshot(contextDir string, files []string) (string, error) {
	stamp := time.Now().UTC().Format(cfgInit.BackupTimestampLayout)
	backupDir := filepath.Join(
		contextDir,
		cfgInit.BackupDirPrefix+stamp,
	)
	if mkErr := ctxIo.SafeMkdirAll(backupDir, fs.PermExec); mkErr != nil {
		return "", errInit.BackupMkdir(backupDir, mkErr)
	}

	for _, name := range files {
		src := filepath.Join(contextDir, name)
		data, readErr := ctxIo.SafeReadUserFile(src)
		if readErr != nil {
			if os.IsNotExist(readErr) {
				continue
			}
			return "", errInit.BackupRead(src, readErr)
		}
		dst := filepath.Join(backupDir, name)
		if writeErr := ctxIo.SafeWriteFile(dst, data, fs.PermFile); writeErr != nil {
			return "", errInit.BackupWrite(dst, writeErr)
		}
	}

	return backupDir, nil
}
