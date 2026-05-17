//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

import (
	"path/filepath"
	"time"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	errKbSC "github.com/ActiveMemory/ctx/internal/err/kb/sourcecoverage"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Advance writes (or updates) a row in the ledger. When the
// source already has a row, the transition is validated; when
// the source is new, the row is appended only when the State
// is `discovered` or `admitted` (the entry points of the state
// machine).
//
// Parameters:
//   - ledgerPath: full path to
//     `.context/kb/source-coverage.md`.
//   - row: row to write. Updated is set to time.Now() when
//     zero.
//
// Returns:
//   - error: [errKbSC.ErrIllegalTransition],
//     [errKbSC.ErrUnknownSource], or wrapped I/O failures.
func Advance(ledgerPath string, row Row) error {
	if row.Updated.IsZero() {
		row.Updated = time.Now().UTC().Truncate(time.Second)
	}
	rows, readErr := Read(ledgerPath)
	if readErr != nil {
		return readErr
	}

	idx := -1
	for i, r := range rows {
		if r.Source == row.Source {
			idx = i
			break
		}
	}
	if idx < 0 {
		switch row.State {
		case cfgKB.StateDiscovered, cfgKB.StateAdmitted:
			rows = append(rows, row)
		default:
			return errKbSC.UnknownSource(row.Source, row.State)
		}
	} else {
		if !ValidTransition(rows[idx].State, row.State) {
			return errKbSC.IllegalTransition(
				rows[idx].State, row.State, row.Source,
			)
		}
		rows[idx] = row
	}

	if mkErr := ctxIo.SafeMkdirAll(
		filepath.Dir(ledgerPath), cfgFs.PermExec,
	); mkErr != nil {
		return errKbSC.MkdirLedgerDir(mkErr)
	}
	rendered := render(rows)
	if writeErr := ctxIo.SafeWriteFile(
		ledgerPath, []byte(rendered), cfgFs.PermSecret,
	); writeErr != nil {
		return errKbSC.WriteLedger(writeErr)
	}
	return nil
}
