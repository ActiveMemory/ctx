//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contradiction

import (
	cfgKbC "github.com/ActiveMemory/ctx/internal/config/kb/contradiction"
	"github.com/ActiveMemory/ctx/internal/entity"
	errKbC "github.com/ActiveMemory/ctx/internal/err/kb/contradiction"
	"github.com/ActiveMemory/ctx/internal/write/kb/row"
)

// Append writes one row to the contradictions artifact at
// path, allocating the next monotonic `C-###` ID. When the
// file does not exist, it is created with the schema header
// and the first row is assigned `C-001`. The write opens the
// file with O_CREATE|O_APPEND|O_WRONLY; idempotency at the
// call-site is the caller's responsibility.
//
// Parameters:
//   - path: absolute path to `.context/kb/contradictions.md`.
//   - r: row content; ID is filled in by Append.
//
// Returns:
//   - string: the allocated `C-###` ID.
//   - error: wrapped I/O / parse failures.
func Append(path string, r Row) (string, error) {
	return row.Append(path, entity.KBRowHooks{
		Header:   cfgKbC.TableHeader,
		NextID:   nextID,
		Render:   func(id string) string { return renderRow(id, r) },
		ErrMkdir: errKbC.MkdirDir,
		ErrRead:  errKbC.ReadFile,
		ErrOpen:  errKbC.OpenFile,
		ErrWrite: errKbC.WriteRow,
	})
}
