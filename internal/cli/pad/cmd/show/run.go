//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run prints the raw text of entry by stable ID.
//
// Parameters:
//   - cmd: Cobra command for output
//   - id: Stable entry ID
//   - outPath: File path for blob output (empty for stdout)
//
// Returns:
//   - error: Non-nil on invalid ID, read or write failure
func Run(cmd *cobra.Command, id int, outPath string) error {
	if _, ctxErr := rc.RequireContextDir(); ctxErr != nil {
		cmd.SilenceUsage = true
		return ctxErr
	}
	entries, readErr := store.ReadEntriesWithIDs()
	if readErr != nil {
		return readErr
	}

	idx := parse.FindByID(entries, id)
	if idx < 0 {
		return errPad.EntryNotFound(id)
	}

	entry := entries[idx].Content

	if _, data, ok := blob.Split(entry); ok {
		if outPath != "" {
			if writeErr := ctxIo.SafeWriteFile(
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
