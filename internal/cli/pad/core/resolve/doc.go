//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve provides scratchpad display
// conversion helpers.
//
// When listing scratchpad entries, raw entry strings
// need to be converted to a human-readable form. Blob
// entries contain base64 payloads that should not be
// shown verbatim. This package maps raw entries to
// their display representations.
//
// # Display Conversion
//
// [DisplayAll] is the sole exported function. It
// accepts a slice of raw entry strings (typically
// from decryption) and returns a new slice of the
// same length where each entry has been processed
// through blob.DisplayEntry. Plain text entries pass
// through unchanged. Blob entries are replaced with
// "label [BLOB]" to indicate binary content without
// dumping encoded data.
//
// # Data Flow
//
// The cmd/pad layer decrypts the scratchpad via the
// store package, calls DisplayAll to prepare entries
// for output, and passes the result to the write/pad
// package for formatted display. DisplayAll is a pure
// transformation with no side effects.
package resolve
