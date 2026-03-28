//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"github.com/spf13/cobra"

	coreEdit "github.com/ActiveMemory/ctx/internal/cli/pad/core/edit"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// RunEdit replaces the entry at 1-based position n with new text.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - text: Replacement text
//
// Returns:
//   - error: Non-nil on invalid index or read/write failure
func RunEdit(cmd *cobra.Command, n int, text string) error {
	entries, editErr := coreEdit.Replace(n, text)
	if editErr != nil {
		return editErr
	}
	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}
	writePad.EntryUpdated(cmd, n)
	return nil
}

// RunEditAppend appends text to the entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - text: Text to append
//
// Returns:
//   - error: Non-nil on invalid index, blob entry, or read/write failure
func RunEditAppend(cmd *cobra.Command, n int, text string) error {
	entries, editErr := coreEdit.Append(n, text)
	if editErr != nil {
		return editErr
	}
	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}
	writePad.EntryUpdated(cmd, n)
	return nil
}

// RunEditPrepend prepends text to the entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - text: Text to prepend
//
// Returns:
//   - error: Non-nil on invalid index, blob entry, or read/write failure
func RunEditPrepend(cmd *cobra.Command, n int, text string) error {
	entries, editErr := coreEdit.Prepend(n, text)
	if editErr != nil {
		return editErr
	}
	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}
	writePad.EntryUpdated(cmd, n)
	return nil
}

// RunEditBlob replaces the file content and/or label of a blob entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - filePath: New file content path (empty to keep existing)
//   - labelText: New label (empty to keep existing)
//
// Returns:
//   - error: Non-nil on invalid index, non-blob entry, or read/write failure
func RunEditBlob(cmd *cobra.Command, n int, filePath, labelText string) error {
	entries, editErr := coreEdit.UpdateBlob(n, filePath, labelText)
	if editErr != nil {
		return editErr
	}
	if writeErr := store.WriteEntries(cmd, entries); writeErr != nil {
		return writeErr
	}
	writePad.EntryUpdated(cmd, n)
	return nil
}
