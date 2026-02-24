//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"encoding/base64"
	"strings"
)

// BlobSep separates the label from the base64-encoded file content.
const BlobSep = ":::"

// MaxBlobSize is the maximum file size (pre-encoding) allowed for blob entries.
const MaxBlobSize = 64 * 1024

// isBlob returns true if the entry contains the blob separator.
func isBlob(entry string) bool {
	return strings.Contains(entry, BlobSep)
}

// splitBlob parses a blob entry into its label and decoded data.
// Returns ok=false for non-blob entries or malformed base64.
func splitBlob(entry string) (label string, data []byte, ok bool) {
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

// makeBlob creates a blob entry string from a label and file data.
func makeBlob(label string, data []byte) string {
	return label + BlobSep + base64.StdEncoding.EncodeToString(data)
}

// displayEntry returns a display-friendly version of an entry.
// Blob entries show "label [BLOB]", plain entries are returned as-is.
func displayEntry(entry string) string {
	if label, _, ok := splitBlob(entry); ok {
		return label + " [BLOB]"
	}
	return entry
}
