//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/err/fs"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/io"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// runAdd appends a new entry and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - text: Entry text to add
//
// Returns:
//   - error: Non-nil on read/write failure
func runAdd(cmd *cobra.Command, text string) error {
	entries, err := store.ReadEntries()
	if err != nil {
		return err
	}

	entries = append(entries, text)

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.EntryAdded(cmd, len(entries))
	return nil
}

// runAddBlob reads a file, encodes it as a blob entry, and appends it.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: Blob label (filename)
//   - filePath: Path to the file to ingest
//
// Returns:
//   - error: Non-nil on read/write failure or file too large
func runAddBlob(cmd *cobra.Command, label, filePath string) error {
	data, err := io.SafeReadUserFile(filePath)
	if err != nil {
		return fs.ReadFile(err)
	}

	if len(data) > pad.MaxBlobSize {
		return errPad.FileTooLarge(len(data), pad.MaxBlobSize)
	}

	entries, readErr := store.ReadEntries()
	if readErr != nil {
		return readErr
	}

	entries = append(entries, blob.MakeBlob(label, data))

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.EntryAdded(cmd, len(entries))
	return nil
}
