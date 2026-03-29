//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

// Scratchpad blob constants.
const (
	// BlobSep separates the label from the base64-encoded file content.
	BlobSep = ":::"
	// MaxBlobSize is the maximum file size (pre-encoding)
	// allowed for blob entries.
	MaxBlobSize = 64 * 1024
	// BlobTag is the display tag appended to blob labels.
	BlobTag = " [BLOB]"
)
