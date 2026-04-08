//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// EntryWithID appends a text entry with a stable ID and
// returns the updated list and the new entry's ID.
//
// Parameters:
//   - text: Entry text to add
//
// Returns:
//   - []parse.Entry: Updated entries
//   - int: ID assigned to the new entry
//   - error: Non-nil on entry load failure
func EntryWithID(
	text string,
) ([]parse.Entry, int, error) {
	entries, loadErr := store.ReadEntriesWithIDs()
	if loadErr != nil {
		return nil, 0, loadErr
	}
	id := parse.NextID(entries)
	entries = append(entries, parse.Entry{
		ID: id, Content: text,
	})
	return entries, id, nil
}

// BlobWithID reads a file, validates size, encodes as blob,
// and returns the updated entries with stable IDs.
//
// Parameters:
//   - label: Blob label (filename)
//   - filePath: Path to the file to ingest
//
// Returns:
//   - []parse.Entry: Updated entries
//   - int: ID assigned to the new entry
//   - error: Non-nil on read failure or file too large
func BlobWithID(
	label, filePath string,
) ([]parse.Entry, int, error) {
	data, readErr := internalIo.SafeReadUserFile(filePath)
	if readErr != nil {
		return nil, 0, errFs.ReadFile(readErr)
	}

	if len(data) > pad.MaxBlobSize {
		return nil, 0, errPad.FileTooLarge(
			len(data), pad.MaxBlobSize)
	}

	entries, loadErr := store.ReadEntriesWithIDs()
	if loadErr != nil {
		return nil, 0, loadErr
	}

	id := parse.NextID(entries)
	entries = append(entries, parse.Entry{
		ID: id, Content: blob.Make(label, data),
	})
	return entries, id, nil
}
