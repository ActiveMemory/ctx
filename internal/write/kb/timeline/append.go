//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package timeline

import (
	"os"
	"path/filepath"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgKbTL "github.com/ActiveMemory/ctx/internal/config/kb/timeline"
	errKbTL "github.com/ActiveMemory/ctx/internal/err/kb/timeline"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Append writes one row to the timeline artifact at path. When
// the file does not exist, it is created with the schema
// header and then the row is appended. The write opens the
// file with O_CREATE|O_APPEND|O_WRONLY; idempotency at the
// call-site is the caller's responsibility.
//
// Parameters:
//   - path: absolute path to `.context/kb/timeline.md`.
//   - row: row content.
//
// Returns:
//   - error: wrapped I/O failures.
func Append(path string, row Row) error {
	if mkErr := ctxIo.SafeMkdirAll(
		filepath.Dir(path), cfgFs.PermExec,
	); mkErr != nil {
		return errKbTL.MkdirDir(mkErr)
	}
	needsHeader := false
	if _, statErr := ctxIo.SafeStat(path); statErr != nil {
		if !os.IsNotExist(statErr) {
			return errKbTL.ReadFile(statErr)
		}
		needsHeader = true
	}

	f, openErr := ctxIo.SafeAppendFile(path, cfgFs.PermSecret)
	if openErr != nil {
		return errKbTL.OpenFile(openErr)
	}
	defer func() { _ = f.Close() }()

	if needsHeader {
		if _, writeErr := f.WriteString(
			cfgKbTL.TableHeader,
		); writeErr != nil {
			return errKbTL.WriteRow(writeErr)
		}
	}
	if _, writeErr := f.WriteString(
		renderRow(row),
	); writeErr != nil {
		return errKbTL.WriteRow(writeErr)
	}
	return nil
}
