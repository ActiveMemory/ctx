//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package export implements the "ctx pad export"
// subcommand for writing blob entries to disk.
//
// # Behavior
//
// The command scans the scratchpad for blob entries
// and writes each blob's binary data to a file in the
// target directory. The target defaults to the current
// directory (".") but accepts an optional positional
// argument for a different path.
//
// When a destination file already exists, the command
// appends a timestamp suffix to avoid overwriting,
// unless --force is set. The --dry-run flag previews
// the export plan without writing any files.
//
// Plain-text entries are silently skipped; only blob
// entries produce output files.
//
// # Flags
//
//	--force, -f    Overwrite existing files instead
//	               of generating timestamped names.
//	--dry-run      Print the export plan without
//	               writing files to disk.
//
// # Output
//
// Each exported blob prints a confirmation line with
// the blob label. At the end, a summary line reports
// the total count of exported (or planned) files.
// In dry-run mode, name collisions are noted with
// the alternative filename that would be used.
//
// # Delegation
//
// Export planning is handled by [coreExport.Plan].
// File I/O uses [ctxIo.SafeWriteFile] with secret
// permissions. Output formatting goes through the
// [writePad] and [writeExport] packages.
package export
