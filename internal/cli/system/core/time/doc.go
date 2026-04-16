//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package time returns the most recent modification time
// across context files in the .context/ directory. The
// persistence nudge system uses this to detect whether
// context files have been updated since the last nudge.
//
// # Modification Time Scan
//
// [GetLatestMtime] scans all .md files in the given
// context directory and returns the Unix timestamp of the
// most recently modified file. Non-markdown files and
// subdirectories are skipped. Returns 0 when the
// directory cannot be read or contains no markdown files.
//
// The function reads directory entries and stats each
// markdown file individually, selecting the maximum
// mtime. This avoids parsing file contents and keeps the
// check lightweight enough to run on every prompt cycle.
//
// # Usage Pattern
//
// The persistence hook compares the latest mtime against
// the stored LastMtime in the persistence state. When
// the mtime has advanced, the hook resets its nudge
// counter because the agent (or user) has already
// updated context files, making a nudge unnecessary.
package time
