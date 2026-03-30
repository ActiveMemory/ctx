//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"os"

	"github.com/spf13/cobra"

	coreImp "github.com/ActiveMemory/ctx/internal/cli/pad/core/imp"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Lines imports pad entries from a line-delimited file or stdin.
//
// Parameters:
//   - cmd: Cobra command for status output
//   - file: Path to the input file, or "-" for stdin
//
// Returns:
//   - error: Non-nil on open, read, or write failure
func Lines(cmd *cobra.Command, file string) error {
	var r *os.File
	if file == cli.StdinSentinel {
		r = os.Stdin
	} else {
		f, openErr := internalIo.SafeOpenUserFile(file)
		if openErr != nil {
			return errFs.OpenFile(file, openErr)
		}
		defer func() {
			if cErr := f.Close(); cErr != nil {
				writePad.ErrImportCloseWarning(cmd, file, cErr)
			}
		}()
		r = f
	}

	entries, count, readErr := coreImp.FromReader(r)
	if readErr != nil {
		return readErr
	}

	if count == 0 {
		writePad.ImportNone(cmd)
		return nil
	}

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.ImportDone(cmd, count)
	return nil
}

// Blobs imports pad entries from files in a directory.
//
// Parameters:
//   - cmd: Cobra command for status output
//   - path: Directory path containing blob files to import
//
// Returns:
//   - error: Non-nil on directory read or entry write failure
func Blobs(cmd *cobra.Command, path string) error {
	entries, added, results, dirErr := coreImp.FromDirectory(path)
	if dirErr != nil {
		return dirErr
	}

	skipped := 0
	for _, r := range results {
		switch {
		case r.Err != nil:
			writePad.ErrImportBlobSkipped(cmd, r.Name, r.Err)
			skipped++
		case r.TooLarge:
			writePad.ErrImportBlobTooLarge(cmd, r.Name, pad.MaxBlobSize)
			skipped++
		case r.Added:
			writePad.ImportBlobAdded(cmd, r.Name)
		}
	}

	if added > 0 {
		if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
			return writeErr
		}
	}

	writePad.ImportBlobSummary(cmd, added, skipped)
	return nil
}
