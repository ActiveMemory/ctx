//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package load orchestrates scratchpad import from
// files, stdin, or directories.
//
// This package sits between the cmd layer and the
// lower-level imp (import) package. It handles file
// opening, coordinates parsing and storage, and
// produces user-facing output via the write/pad
// package.
//
// # Line-Delimited Import
//
// [Lines] imports pad entries from a line-delimited
// file or stdin. When the file argument is the stdin
// sentinel ("-"), it reads from os.Stdin; otherwise it
// opens the named file via SafeOpenUserFile and defers
// close. It delegates parsing to imp.FromReader, which
// returns the updated entries and count. If no entries
// were added, it reports "none imported" and returns.
// Otherwise it writes the entries via store.WriteEntries
// and reports the count.
//
// # Directory Blob Import
//
// [Blobs] imports pad entries from files in a
// directory. It delegates to imp.FromDirectory, which
// reads first-level regular files as blob entries. For
// each BlobResult, Blobs reports errors, size limit
// violations, and successful additions via the
// write/pad package. If any blobs were added, it
// writes the combined entries to storage.
//
// # Data Flow
//
// Both functions follow the same pattern: open input,
// delegate to imp for parsing, write via store, report
// via write/pad. The cmd/pad layer calls Lines or
// Blobs based on user flags.
package load
