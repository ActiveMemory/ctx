//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

import (
	cfgKbDD "github.com/ActiveMemory/ctx/internal/config/kb/decision"
	"github.com/ActiveMemory/ctx/internal/entity"
	errKbDD "github.com/ActiveMemory/ctx/internal/err/kb/decision"
	"github.com/ActiveMemory/ctx/internal/write/kb/row"
)

// Append writes one row to the domain-decisions artifact at
// path, allocating the next monotonic `DD-###` ID. When the
// file does not exist, it is created with the schema header
// and the first row is assigned `DD-001`. The write opens the
// file with O_CREATE|O_APPEND|O_WRONLY; idempotency at the
// call-site is the caller's responsibility.
//
// Parameters:
//   - path: absolute path to
//     `.context/kb/domain-decisions.md`.
//   - r: row content; ID is filled in by Append.
//
// Returns:
//   - string: the allocated `DD-###` ID.
//   - error: wrapped I/O / parse failures.
func Append(path string, r Row) (string, error) {
	return row.Append(path, entity.KBRowHooks{
		Header:   cfgKbDD.TableHeader,
		NextID:   nextID,
		Render:   func(id string) string { return renderRow(id, r) },
		ErrMkdir: errKbDD.MkdirDir,
		ErrRead:  errKbDD.ReadFile,
		ErrOpen:  errKbDD.OpenFile,
		ErrWrite: errKbDD.WriteRow,
	})
}
