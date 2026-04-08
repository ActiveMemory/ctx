//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	coreAdd "github.com/ActiveMemory/ctx/internal/cli/pad/core/add"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run appends a new text or blob entry to the scratchpad.
//
// When filePath is non-empty, the entry is imported as a blob
// with text as the label. Otherwise text is added as plain.
//
// Parameters:
//   - cmd: Cobra command for output
//   - text: Entry text or blob label
//   - filePath: Blob file path (empty for plain text)
//
// Returns:
//   - error: Non-nil on read/write failure or too large
func Run(cmd *cobra.Command, text, filePath string) error {
	var entries []parse.Entry
	var id int
	var addErr error

	if filePath != "" {
		entries, id, addErr = coreAdd.BlobWithID(
			text, filePath)
	} else {
		entries, id, addErr = coreAdd.EntryWithID(text)
	}
	if addErr != nil {
		return addErr
	}

	writeErr := store.WriteEntriesWithIDs(cmd, entries)
	if writeErr != nil {
		return writeErr
	}

	writePad.EntryAdded(cmd, id)
	return nil
}
