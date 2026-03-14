//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
)

// runEdit replaces entry at 1-based position n with new text.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - text: Replacement text
//
// Returns:
//   - error: Non-nil on invalid index or read/write failure
func runEdit(cmd *cobra.Command, n int, text string) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if validErr := core.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}

	entries[n-1] = text

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	write.PadEntryUpdated(cmd, n)
	return nil
}

// runEditAppend appends text to the entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - text: Text to append
//
// Returns:
//   - error: Non-nil on invalid index, blob entry, or read/write failure
func runEditAppend(cmd *cobra.Command, n int, text string) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if validErr := core.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}

	if core.IsBlob(entries[n-1]) {
		return ctxerr.BlobAppendNotAllowed()
	}

	entries[n-1] = entries[n-1] + " " + text

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	write.PadEntryUpdated(cmd, n)
	return nil
}

// runEditPrepend prepends text to the entry at 1-based position n.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - text: Text to prepend
//
// Returns:
//   - error: Non-nil on invalid index, blob entry, or read/write failure
func runEditPrepend(cmd *cobra.Command, n int, text string) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if validErr := core.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}

	if core.IsBlob(entries[n-1]) {
		return ctxerr.BlobPrependNotAllowed()
	}

	entries[n-1] = text + " " + entries[n-1]

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	write.PadEntryUpdated(cmd, n)
	return nil
}

// runEditBlob replaces the file content and/or label of a blob entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - n: 1-based entry index
//   - filePath: New file content path (empty to keep existing)
//   - labelText: New label (empty to keep existing)
//
// Returns:
//   - error: Non-nil on invalid index, non-blob entry, or read/write failure
func runEditBlob(cmd *cobra.Command, n int, filePath, labelText string) error {
	entries, err := core.ReadEntries()
	if err != nil {
		return err
	}

	if validErr := core.ValidateIndex(n, entries); validErr != nil {
		return validErr
	}

	oldLabel, oldData, ok := core.SplitBlob(entries[n-1])
	if !ok {
		return ctxerr.NotBlobEntry(n)
	}

	newLabel := oldLabel
	newData := oldData

	if labelText != "" {
		newLabel = labelText
	}

	if filePath != "" {
		data, readErr := io.SafeReadUserFile(filePath)
		if readErr != nil {
			return ctxerr.ReadFile(readErr)
		}
		if len(data) > pad.MaxBlobSize {
			return ctxerr.FileTooLarge(len(data), pad.MaxBlobSize)
		}
		newData = data
	}

	entries[n-1] = core.MakeBlob(newLabel, newData)

	if writeErr := core.WriteEntries(entries); writeErr != nil {
		return writeErr
	}

	write.PadEntryUpdated(cmd, n)
	return nil
}
