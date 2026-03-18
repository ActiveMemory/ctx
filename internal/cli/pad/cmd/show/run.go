//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	fs2 "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/write/pad"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
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
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		return ctxerr.EntryRange(n, 0)
	}

	if validErr := core.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}

	entry := entries[n-1]

	if _, data, ok := core.SplitBlob(entry); ok {
		if outPath != "" {
			if writeErr := os.WriteFile(
				outPath, data, fs.PermSecret,
			); writeErr != nil {
				return fs2.WriteFileFailed(writeErr)
			}
			pad.PadBlobWritten(cmd, len(data), outPath)
			return nil
		}
		cmd.Print(string(data))
		return nil
	}

	if outPath != "" {
		return ctxerr.OutFlagRequiresBlob()
	}

	cmd.Println(entry)
	return nil
}
