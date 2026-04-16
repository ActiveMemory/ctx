//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add implements the "ctx pad add" subcommand
// for appending entries to the encrypted scratchpad.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// the text to store. When run without flags, it creates
// a plain-text entry. When the --file (-f) flag is
// provided, the text argument becomes a label and the
// file contents are stored as a binary blob.
//
// Each new entry receives a stable auto-incrementing
// ID that persists across additions and deletions.
// The command prints a confirmation with the assigned
// ID on success.
//
// # Flags
//
//	--file, -f <path>   Import a file as a blob entry;
//	                     the positional arg becomes the
//	                     label for the blob.
//
// # Output
//
// On success, prints a one-line confirmation showing
// the newly assigned entry ID. On failure, returns an
// error for oversized content or read/write problems.
//
// # Delegation
//
// Entry creation is handled by [coreAdd.EntryWithID]
// and [coreAdd.BlobWithID] in the core/add package.
// Persistence goes through [store.WriteEntriesWithIDs].
// User-facing output is routed through [writePad].
package add
