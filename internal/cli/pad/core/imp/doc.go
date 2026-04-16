//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package imp provides import logic for the
// scratchpad.
//
// The "ctx pad import" and "ctx pad load" commands
// ingest entries from external sources into the
// scratchpad. This package reads input, validates it,
// and returns updated entry slices without writing.
//
// # Line-Based Import
//
// [FromReader] reads non-empty lines from an
// io.Reader (file or stdin), trims whitespace, appends
// each line as a new text entry to the existing
// scratchpad entries loaded via store.ReadEntries, and
// returns the combined slice with the count of new
// entries. Empty lines are silently skipped.
//
// # Directory-Based Blob Import
//
// [FromDirectory] reads first-level regular files from
// a directory path. For each file it:
//
//  1. Reads the content via SafeReadFile.
//  2. Checks the size against pad.MaxBlobSize; files
//     that exceed the limit are recorded as TooLarge.
//  3. Encodes passing files as blob entries via
//     blob.Make with the filename as the label.
//
// It returns the updated entries, the count of added
// blobs, and a slice of [BlobResult] values for
// per-file reporting. Sub-directories and non-regular
// files are silently skipped.
//
// # BlobResult Type
//
// [BlobResult] carries per-file outcomes: Name (source
// filename), Err (read error), TooLarge (size limit
// exceeded), and Added (successfully imported).
package imp
