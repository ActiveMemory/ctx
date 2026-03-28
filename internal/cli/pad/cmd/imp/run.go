//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package imp

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

// RunImport reads lines from a file (or stdin) and appends them as entries.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: File path or "-" for stdin
//
// Returns:
//   - error: Non-nil on read/write failure
func RunImport(cmd *cobra.Command, file string) error {
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

// RunImportBlobs reads first-level files from a directory and imports
// each as a blob entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: Directory path containing files to import
//
// Returns:
//   - error: Non-nil on read/write failure
func RunImportBlobs(cmd *cobra.Command, path string) error {
	entries, added, results, dirErr := coreImp.FromDirectory(path)
	if dirErr != nil {
		return dirErr
	}

	// Report per-file outcomes.
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
