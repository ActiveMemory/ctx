//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package row

import (
	"os"
	"path/filepath"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/entity"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Append performs the shared append flow for a monotonic-ID
// kb tabular artifact (contradictions, domain-decisions,
// outstanding-questions). Returns the allocated ID.
//
// Parameters:
//   - path: absolute path to the artifact file.
//   - h: per-artifact hooks; see [entity.KBRowHooks].
//
// Returns:
//   - string: the allocated ID for this row.
//   - error: wrapped via the matching Err* constructor in h.
func Append(path string, h entity.KBRowHooks) (_ string, err error) {
	if mkErr := ctxIo.SafeMkdirAll(
		filepath.Dir(path), cfgFs.PermExec,
	); mkErr != nil {
		return "", h.ErrMkdir(mkErr)
	}

	needsHeader := false
	existing, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		if !os.IsNotExist(readErr) {
			return "", h.ErrRead(readErr)
		}
		needsHeader = true
	}
	id, idErr := h.NextID(existing)
	if idErr != nil {
		return "", idErr
	}

	f, openErr := ctxIo.SafeAppendFile(path, cfgFs.PermSecret)
	if openErr != nil {
		return "", h.ErrOpen(openErr)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = h.ErrWrite(cerr)
		}
	}()

	if needsHeader {
		if _, wErr := f.WriteString(h.Header); wErr != nil {
			return "", h.ErrWrite(wErr)
		}
	}
	if _, wErr := f.WriteString(h.Render(id)); wErr != nil {
		return "", h.ErrWrite(wErr)
	}
	return id, nil
}
