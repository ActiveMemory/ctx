//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sourcecoverage

import (
	"errors"
	"os"

	errKbSC "github.com/ActiveMemory/ctx/internal/err/kb/sourcecoverage"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Read parses the ledger at ledgerPath into a Row slice.
// Returns an empty slice (no error) when the file does not
// exist.
//
// Parameters:
//   - ledgerPath: full path to
//     `.context/kb/source-coverage.md`.
//
// Returns:
//   - []Row: parsed ledger rows in source-order.
//   - error: wrapped [errKbSC.ReadLedger] on I/O failure.
func Read(ledgerPath string) ([]Row, error) {
	raw, readErr := ctxIo.SafeReadUserFile(ledgerPath)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errKbSC.ReadLedger(readErr)
	}
	return parse(string(raw)), nil
}
