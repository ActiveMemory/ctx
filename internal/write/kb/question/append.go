//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package question

import (
	cfgKbQ "github.com/ActiveMemory/ctx/internal/config/kb/question"
	"github.com/ActiveMemory/ctx/internal/entity"
	errKbQ "github.com/ActiveMemory/ctx/internal/err/kb/question"
	"github.com/ActiveMemory/ctx/internal/write/kb/row"
)

// Append writes one row to the outstanding-questions artifact
// at path, allocating the next monotonic `Q-###` ID. When the
// file does not exist, it is created with the schema header
// and the first row is assigned `Q-001`. The write opens the
// file with O_CREATE|O_APPEND|O_WRONLY; idempotency at the
// call-site is the caller's responsibility.
//
// Parameters:
//   - path: absolute path to
//     `.context/kb/outstanding-questions.md`.
//   - r: row content; ID is filled in by Append.
//
// Returns:
//   - string: the allocated `Q-###` ID.
//   - error: wrapped I/O / parse failures.
func Append(path string, r Row) (string, error) {
	return row.Append(path, entity.KBRowHooks{
		Header:   cfgKbQ.TableHeader,
		NextID:   nextID,
		Render:   func(id string) string { return renderRow(id, r) },
		ErrMkdir: errKbQ.MkdirDir,
		ErrRead:  errKbQ.ReadFile,
		ErrOpen:  errKbQ.OpenFile,
		ErrWrite: errKbQ.WriteRow,
	})
}
