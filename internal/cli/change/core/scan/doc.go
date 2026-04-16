//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package scan queries the filesystem and git history
// for changes since a reference time.
//
// # FindContextChanges
//
// [FindContextChanges] reads the .context/ directory
// and returns Markdown files whose modification time is
// after the given reference time. Results are sorted by
// modification time descending (most recent first).
// Directories and non-Markdown files are skipped.
//
// # SummarizeCodeChanges
//
// [SummarizeCodeChanges] produces a CodeSummary by
// running git log commands. It collects four pieces of
// information:
//
//   - Commit count since the reference time.
//   - Latest commit message (first line of oneline).
//   - Unique top-level directories touched by commits.
//   - Unique author names from the commit history.
//
// All git failures produce an empty summary rather than
// an error, so the change command works gracefully in
// non-git directories.
//
// # Helper Functions
//
// [GitLogSince] wraps execGit.LogSince with a --since
// filter derived from the reference time. The time is
// formatted as RFC 3339 internally so no caller input
// reaches exec.Command, satisfying gosec G204.
//
// [UniqueTopDirs] extracts unique top-level directory
// names from newline-separated file paths. It splits
// each path at the first "/" and deduplicates.
//
// [UniqueLines] returns sorted unique non-empty lines
// from newline-separated output, used for deduplicating
// author names.
//
// # Data Flow
//
// The cmd/ layer calls FindContextChanges and
// SummarizeCodeChanges with the resolved reference
// time, then passes both results to the render
// subpackage for formatting.
package scan
