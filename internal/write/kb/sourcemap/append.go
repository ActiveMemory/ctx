//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcemap

import (
	"os"
	"path/filepath"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgKbSM "github.com/ActiveMemory/ctx/internal/config/kb/sourcemap"
	errKbSM "github.com/ActiveMemory/ctx/internal/err/kb/sourcemap"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Append writes one row to the source-map artifact at path.
// When the file does not exist, it is created with the schema
// header and then the row is appended. The write opens the
// file with O_CREATE|O_APPEND|O_WRONLY; idempotency at the
// call-site is the caller's responsibility.
//
// Parameters:
//   - path: absolute path to `.context/kb/source-map.md`.
//   - row: row content.
//
// Returns:
//   - error: wrapped I/O failures.
func Append(path string, row Row) (err error) {
	if mkErr := ctxIo.SafeMkdirAll(
		filepath.Dir(path), cfgFs.PermExec,
	); mkErr != nil {
		return errKbSM.MkdirDir(mkErr)
	}
	needsHeader := false
	if _, statErr := ctxIo.SafeStat(path); statErr != nil {
		if !os.IsNotExist(statErr) {
			return errKbSM.ReadFile(statErr)
		}
		needsHeader = true
	}

	f, openErr := ctxIo.SafeAppendFile(path, cfgFs.PermSecret)
	if openErr != nil {
		return errKbSM.OpenFile(openErr)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = errKbSM.WriteRow(cerr)
		}
	}()

	if needsHeader {
		if _, writeErr := f.WriteString(
			cfgKbSM.TableHeader,
		); writeErr != nil {
			return errKbSM.WriteRow(writeErr)
		}
	}
	if _, writeErr := f.WriteString(
		renderRow(row),
	); writeErr != nil {
		return errKbSM.WriteRow(writeErr)
	}
	return nil
}
