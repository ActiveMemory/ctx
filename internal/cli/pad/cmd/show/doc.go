//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package show implements the "ctx pad show" subcommand
// for displaying a single scratchpad entry by stable
// ID.
//
// # Behavior
//
// The command accepts exactly one positional argument:
// the stable entry ID (1-based integer). It looks up
// the entry and prints its raw content to stdout
// without any numbering prefix, making the output
// suitable for piping into other commands:
//
//	ctx pad edit 1 --append "$(ctx pad show 3)"
//
// For blob entries, the binary data is written to
// stdout by default. When the --out flag is set, the
// blob data is written to the specified file path
// instead. The --out flag is only valid for blob
// entries; using it with a plain-text entry returns
// an error.
//
// # Flags
//
//	--out <path>   Write blob data to a file instead
//	               of stdout. Only valid for blob
//	               entries.
//
// # Output
//
// For plain-text entries, prints the entry content
// on stdout. For blobs without --out, writes raw
// binary data to stdout. For blobs with --out, prints
// a confirmation showing byte count and file path.
//
// # Delegation
//
// Entry lookup uses [parse.FindByID]. Blob detection
// and splitting are handled by [blob.Split]. File
// output uses [ctxIo.SafeWriteFile] with secret
// permissions. Display is routed through [pad].
package show
