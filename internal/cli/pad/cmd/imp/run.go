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

// Run imports entries into the scratchpad from a file, stdin, or directory.
//
// When blobs is true, imports directory contents as blob entries.
// Otherwise reads lines from a file (or stdin when path is "-").
//
// Parameters:
//   - cmd: Cobra command for output
//   - path: File path, "-" for stdin, or directory path (when blobs is true)
//   - blobs: When true, import directory contents as blob entries
//
// Returns:
//   - error: Non-nil on read/write failure
func Run(cmd *cobra.Command, path string, blobs bool) error {
	if blobs {
		return runBlobs(cmd, path)
	}
	return runLines(cmd, path)
}

func runLines(cmd *cobra.Command, file string) error {
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

func runBlobs(cmd *cobra.Command, path string) error {
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
