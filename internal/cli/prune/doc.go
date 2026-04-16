//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prune implements the ctx prune top-level
// command.
//
// Removes stale per-session state files under
// .context/state/ that have not been touched within the
// configured retention window (default 7 days). Only
// files whose names match UUID patterns are pruned;
// global state files are preserved.
//
// # How It Works
//
// The [Run] function scans the state directory, checks
// each UUID-named file's modification time against the
// cutoff, and removes files older than --days. A
// --dry-run flag previews what would be pruned without
// deleting anything.
//
// # Flags
//
//   - --days: retention window in days (default 7)
//   - --dry-run: preview without deleting
//
// [Cmd] returns the cobra command with --days and
// --dry-run flags. [Run] scans the state directory,
// compares each UUID file's mtime against the cutoff,
// and removes expired entries.
package prune
