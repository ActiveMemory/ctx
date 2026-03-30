//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import "time"

// ContextChange represents a modified context file.
//
// Fields:
//   - Name: Context file name (e.g. "TASKS.md")
//   - ModTime: Last modification timestamp
type ContextChange struct {
	Name    string
	ModTime time.Time
}

// CodeSummary summarizes code changes since the reference time.
//
// Fields:
//   - CommitCount: Number of commits since reference
//   - LatestMsg: Most recent commit message
//   - Dirs: Unique directories with changes
//   - Authors: Unique commit authors
type CodeSummary struct {
	CommitCount int
	LatestMsg   string
	Dirs        []string
	Authors     []string
}
