//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

import (
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// Replace replaces the entry at 1-based position n with new text.
// Returns the updated entries list. The caller owns writing.
//
// Parameters:
//   - n: 1-based entry index
//   - text: Replacement text
//
// Returns:
//   - []string: Updated entries list
//   - error: Non-nil on invalid index or load failure
func Replace(n int, text string) ([]string, error) {
	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, loadErr
	}
	if validErr := validate.Index(n, entries); validErr != nil {
		return nil, validErr
	}
	entries[n-1] = text
	return entries, nil
}

// Append appends text to the entry at 1-based position n.
// For blob entries, the text is appended to the label.
// Returns the updated entries list. The caller owns writing.
//
// Parameters:
//   - n: 1-based entry index
//   - text: Text to append
//
// Returns:
//   - []string: Updated entries list
//   - error: Non-nil on invalid index or load failure
func Append(n int, text string) ([]string, error) {
	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, loadErr
	}
	if validErr := validate.Index(n, entries); validErr != nil {
		return nil, validErr
	}
	if label, data, ok := blob.Split(entries[n-1]); ok {
		entries[n-1] = blob.Make(label+" "+text, data)
	} else {
		entries[n-1] = entries[n-1] + " " + text
	}
	return entries, nil
}

// Prepend prepends text to the entry at 1-based position n.
// For blob entries, the text is prepended to the label.
// Returns the updated entries list. The caller owns writing.
//
// Parameters:
//   - n: 1-based entry index
//   - text: Text to prepend
//
// Returns:
//   - []string: Updated entries list
//   - error: Non-nil on invalid index or load failure
func Prepend(n int, text string) ([]string, error) {
	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, loadErr
	}
	if validErr := validate.Index(n, entries); validErr != nil {
		return nil, validErr
	}
	if label, data, ok := blob.Split(entries[n-1]); ok {
		entries[n-1] = blob.Make(text+" "+label, data)
	} else {
		entries[n-1] = text + " " + entries[n-1]
	}
	return entries, nil
}

// UpdateBlob replaces the file content and/or label of a blob entry.
// Returns the updated entries list. The caller owns writing.
//
// Parameters:
//   - n: 1-based entry index
//   - filePath: New file content path (empty to keep existing)
//   - labelText: New label (empty to keep existing)
//
// Returns:
//   - []string: Updated entries list
//   - error: Non-nil on invalid index, non-blob entry, or read failure
func UpdateBlob(n int, filePath, labelText string) ([]string, error) {
	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, loadErr
	}
	if validErr := validate.Index(n, entries); validErr != nil {
		return nil, validErr
	}

	oldLabel, oldData, ok := blob.Split(entries[n-1])
	if !ok {
		return nil, errPad.NotBlobEntry(n)
	}

	newLabel := oldLabel
	newData := oldData

	if labelText != "" {
		newLabel = labelText
	}

	if filePath != "" {
		data, readErr := internalIo.SafeReadUserFile(filePath)
		if readErr != nil {
			return nil, readErr
		}
		if len(data) > pad.MaxBlobSize {
			return nil, errPad.FileTooLarge(len(data), pad.MaxBlobSize)
		}
		newData = data
	}

	entries[n-1] = blob.Make(newLabel, newData)
	return entries, nil
}
