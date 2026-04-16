//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup provides terminal output for the backup
// command (ctx backup).
//
// # Exported Functions
//
// [ResultLine] prints a single backup result line that
// includes the scope label (e.g. "project" or "global"),
// the archive file path, and the human-readable file
// size. When an SMB destination is configured, the line
// is extended with the remote copy path.
//
// [SkipEntry] writes a notice to an io.Writer when an
// optional archive entry is skipped because its source
// file does not exist on disk. This function accepts an
// io.Writer instead of *cobra.Command because it runs
// during archive assembly before command output is
// available.
//
// # Message Categories
//
//   - Info: backup result with scope, path, and size
//   - Skip: notice when an optional file is missing
//
// # Usage
//
//	backup.ResultLine(cmd, "project", archivePath,
//	    fileSize, smbDest)
//	backup.SkipEntry(w, "optional/file")
package backup
