//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"encoding/base64"
	"strings"
)

// BlobSep separates the label from the base64-encoded file content.
const BlobSep = ":::"

// MaxBlobSize is the maximum file size (pre-encoding) allowed for blob entries.
const MaxBlobSize = 64 * 1024

// IsBlob returns true if the entry contains the blob separator.
//
// Parameters:
//   - entry: Scratchpad entry string
//
// Returns:
//   - bool: True if entry is a blob
func IsBlob(entry string) bool {
	return strings.Contains(entry, BlobSep)
}

// SplitBlob parses a blob entry into its label and decoded data.
//
// Parameters:
//   - entry: Scratchpad entry string
//
// Returns:
//   - label: Blob label (filename)
//   - data: Decoded file content
//   - ok: False for non-blob entries or malformed base64
func SplitBlob(entry string) (label string, data []byte, ok bool) {
	idx := strings.Index(entry, BlobSep)
	if idx < 0 {
		return "", nil, false
	}

	label = entry[:idx]
	encoded := entry[idx+len(BlobSep):]

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", nil, false
	}

	return label, data, true
}

// MakeBlob creates a blob entry string from a label and file data.
//
// Parameters:
//   - label: Blob label (filename)
//   - data: Raw file content
//
// Returns:
//   - string: Formatted blob entry
func MakeBlob(label string, data []byte) string {
	return label + BlobSep + base64.StdEncoding.EncodeToString(data)
}

// DisplayEntry returns a display-friendly version of an entry.
//
// Blob entries show "label [BLOB]", plain entries are returned as-is.
//
// Parameters:
//   - entry: Scratchpad entry string
//
// Returns:
//   - string: Human-readable entry representation
func DisplayEntry(entry string) string {
	if label, _, ok := SplitBlob(entry); ok {
		return label + " [BLOB]"
	}
	return entry
}
