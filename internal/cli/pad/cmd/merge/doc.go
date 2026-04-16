//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package merge implements the "ctx pad merge"
// subcommand for importing entries from external
// scratchpad files into the current pad.
//
// # Behavior
//
// The command reads one or more input files, each
// containing scratchpad entries (optionally encrypted
// with a different key). It deduplicates incoming
// entries against the current pad and appends only
// new ones. Duplicate entries are reported but not
// added.
//
// When the input contains binary blob entries, the
// command warns the user. If a blob label conflicts
// with an existing blob (same label, different data),
// a conflict notice is printed.
//
// # Flags
//
//	--key, -k <path>   Path to the encryption key
//	                    for the input files. When
//	                    omitted, uses the project key.
//	--dry-run           Print the merge summary
//	                    without writing changes.
//
// # Output
//
// Each new entry prints a confirmation line with the
// source file. Duplicates are reported individually.
// A final summary line shows the count of entries
// added and duplicates skipped, plus whether it was
// a dry run.
//
// # Delegation
//
// Key loading and file parsing are handled by the
// core/merge package. Deduplication uses an in-memory
// set built from the current pad. Blob conflict
// detection is provided by [merge.HasBlobConflict].
// Persistence goes through [store.WriteEntries].
package merge
