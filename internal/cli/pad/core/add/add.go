//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// Entry appends a text entry to the scratchpad and returns the
// updated list. The caller owns writing and output.
//
// Parameters:
//   - text: Entry text to add
//
// Returns:
//   - []string: Updated entries list
//   - error: Non-nil on entry load failure
func Entry(text string) ([]string, error) {
	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, loadErr
	}
	return append(entries, text), nil
}

// Blob reads a file, validates its size, encodes it as a blob entry,
// and returns the updated entries list. The caller owns writing and output.
//
// Parameters:
//   - label: Blob label (filename)
//   - filePath: Path to the file to ingest
//
// Returns:
//   - []string: Updated entries list
//   - error: Non-nil on read failure or file too large
func Blob(label, filePath string) ([]string, error) {
	data, readErr := internalIo.SafeReadUserFile(filePath)
	if readErr != nil {
		return nil, errFs.ReadFile(readErr)
	}

	if len(data) > pad.MaxBlobSize {
		return nil, errPad.FileTooLarge(len(data), pad.MaxBlobSize)
	}

	entries, loadErr := store.ReadEntries()
	if loadErr != nil {
		return nil, loadErr
	}

	return append(entries, blob.Make(label, data)), nil
}
