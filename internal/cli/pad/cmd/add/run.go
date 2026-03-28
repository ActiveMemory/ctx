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

// Run appends a new text or blob entry to the scratchpad.
//
// When filePath is non-empty, the entry is imported as a blob with text
// as the label. Otherwise text is added as a plain entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - text: Entry text or blob label
//   - filePath: Blob file path (empty for plain text entry)
//
// Returns:
//   - error: Non-nil on read/write failure or file too large
func Run(cmd *cobra.Command, text, filePath string) error {
	var entries []string
	var addErr error

	if filePath != "" {
		entries, addErr = coreAdd.Blob(text, filePath)
	} else {
		entries, addErr = coreAdd.Entry(text)
	}
	if addErr != nil {
		return addErr
	}

	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}

	writePad.EntryAdded(cmd, len(entries))
	return nil
}
