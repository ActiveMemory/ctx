//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	coreAdd "github.com/ActiveMemory/ctx/internal/cli/pad/core/add"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// RunAdd appends a new entry and prints confirmation.
//
// Parameters:
//   - cmd: Cobra command for output
//   - text: Entry text to add
//
// Returns:
//   - error: Non-nil on read/write failure
func RunAdd(cmd *cobra.Command, text string) error {
	entries, addErr := coreAdd.Entry(text)
	if addErr != nil {
		return addErr
	}

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.EntryAdded(cmd, len(entries))
	return nil
}

// RunAddBlob reads a file, encodes it as a blob entry, and appends it.
//
// Parameters:
//   - cmd: Cobra command for output
//   - label: Blob label (filename)
//   - filePath: Path to the file to ingest
//
// Returns:
//   - error: Non-nil on read/write failure or file too large
func RunAddBlob(cmd *cobra.Command, label, filePath string) error {
	entries, blobErr := coreAdd.Blob(label, filePath)
	if blobErr != nil {
		return blobErr
	}

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.EntryAdded(cmd, len(entries))
	return nil
}
