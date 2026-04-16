//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx pad import"
// subcommand for bulk-loading entries into the
// scratchpad from files, stdin, or directories.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// a file path, the literal "-" for stdin, or a
// directory path when used with the --blob flag.
//
// Without --blob, the command reads the source as a
// text file (or stdin stream) and imports each line
// as a separate plain-text entry. With --blob, it
// treats the argument as a directory and imports each
// file in that directory as a blob entry, using the
// filename as the blob label.
//
// # Flags
//
//	--blob    Import directory contents as blob
//	          entries instead of reading lines from
//	          a text file.
//
// # Output
//
// Delegates all output to the core/load package,
// which prints per-entry confirmations and a summary
// count. Errors from file I/O or oversized entries
// are returned to the caller.
//
// # Delegation
//
// Line-based imports are handled by [load.Lines].
// Blob directory imports use [load.Blobs]. Both
// functions handle reading, entry creation, and
// persistence internally.
package root
