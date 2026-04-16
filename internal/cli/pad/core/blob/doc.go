//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package blob handles binary blob encoding and
// decoding within scratchpad entries.
//
// Scratchpad entries can hold either plain text or
// binary file content. Binary content is stored as a
// base64-encoded string prefixed by a label and a blob
// separator token. This package provides the codec for
// that format.
//
// # Blob Format
//
// A blob entry is a single string in the form:
//
//	<label><BlobSep><base64_data>
//
// The separator token (pad.BlobSep) delimits the
// human-readable label from the encoded payload.
//
// # Detection
//
// [Contains] checks whether an entry string contains
// the blob separator, returning true for blob entries.
// This is a fast check used before attempting the
// more expensive Split operation.
//
// # Parsing
//
// [Split] parses a blob entry into its label and
// decoded byte data. It locates the separator, extracts
// the label prefix, and base64-decodes the suffix. If
// the entry is not a blob or the base64 is malformed,
// it returns ok=false.
//
// # Construction
//
// [Make] creates a blob entry string from a label and
// raw file bytes by concatenating the label, separator,
// and base64-encoded data.
//
// # Display
//
// [DisplayEntry] returns a human-readable form of an
// entry. Blob entries are shown as "label [BLOB]";
// plain text entries are returned unchanged. The
// resolve and list commands use this for output.
package blob
