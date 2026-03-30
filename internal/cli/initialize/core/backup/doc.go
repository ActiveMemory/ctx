//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package backup handles file backup during initialization.
//
// [File] creates a timestamped backup of an existing file before
// overwriting, preserving the user's changes for recovery. Backup
// files use the format name.timestamp.bak and are written alongside
// the original.
package backup
