//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run prints the raw text of entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - outPath: File path for blob output (empty for stdout)
//
// Returns:
//   - error: Non-nil on invalid index, read failure, or write failure
func Run(cmd *cobra.Command, n int, outPath string) error {
	entries, err := store.ReadEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return errPad.EntryRange(n, 0)
	}

	if validErr := validate.Index(n, entries); validErr != nil {
		return validErr
	}

	entry := entries[n-1]

	if _, data, ok := blob.Split(entry); ok {
		if outPath != "" {
			if writeErr := os.WriteFile(
				outPath, data, fs.PermSecret,
			); writeErr != nil {
				return errFs.WriteFileFailed(writeErr)
			}
			pad.BlobWritten(cmd, len(data), outPath)
			return nil
		}
		pad.BlobShow(cmd, data)
		return nil
	}

	if outPath != "" {
		return errPad.OutFlagRequiresBlob()
	}

	pad.EntryShow(cmd, entry)
	return nil
}
