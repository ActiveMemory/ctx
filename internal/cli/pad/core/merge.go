//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"unicode/utf8"

	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/validation"
)

// ReadFileEntries reads a scratchpad file, attempting decryption first.
//
// Parameters:
//   - path: path to the scratchpad file.
//   - key: encryption key (nil to skip the decryption attempt).
//
// Returns:
//   - []string: parsed entries.
//   - error: non-nil if the file cannot be read.
func ReadFileEntries(path string, key []byte) ([]string, error) {
	data, readErr := validation.ReadUserFile(path)
	if readErr != nil {
		return nil, readErr
	}

	if len(data) == 0 {
		return nil, nil
	}

	if key != nil {
		plaintext, decErr := crypto.Decrypt(key, data)
		if decErr == nil {
			return ParseEntries(plaintext), nil
		}
	}

	return ParseEntries(data), nil
}

// LoadMergeKey loads the encryption key for merge input decryption.
//
// Parameters:
//   - keyFile: explicit key file path (empty string = use project key).
//
// Returns:
//   - []byte: the loaded key, or nil if no key is available.
func LoadMergeKey(keyFile string) []byte {
	path := keyFile
	if path == "" {
		path = KeyPath()
	}

	key, loadErr := crypto.LoadKey(path)
	if loadErr != nil {
		return nil
	}
	return key
}

// BuildBlobLabelMap creates a map of blob labels to their full entry strings.
//
// Parameters:
//   - entries: scratchpad entries to scan.
//
// Returns:
//   - map[string]string: blob label to full entry string.
func BuildBlobLabelMap(entries []string) map[string]string {
	labels := make(map[string]string)
	for _, entry := range entries {
		if label, _, ok := SplitBlob(entry); ok {
			labels[label] = entry
		}
	}
	return labels
}

// HasBlobConflict checks if a blob entry has the same label as an existing
// blob but different content. Updates the label map with the new entry.
//
// Parameters:
//   - entry: the new entry to check.
//   - blobLabels: map of existing blob labels to their full entry strings.
//
// Returns:
//   - bool: true if a conflict was detected.
//   - string: the conflicting label (empty if no conflict).
func HasBlobConflict(entry string, blobLabels map[string]string) (bool, string) {
	label, _, ok := SplitBlob(entry)
	if !ok {
		return false, ""
	}

	existing, found := blobLabels[label]
	conflict := found && existing != entry
	blobLabels[label] = entry
	return conflict, label
}

// HasBinaryEntries checks if any entries contain non-UTF-8 bytes.
//
// Parameters:
//   - entries: the parsed entries to check.
//
// Returns:
//   - bool: true if any entry contains non-UTF-8 data.
func HasBinaryEntries(entries []string) bool {
	for _, entry := range entries {
		if !utf8.ValidString(entry) {
			return true
		}
	}
	return false
}
