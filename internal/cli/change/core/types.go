//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import "time"

// ContextChange represents a modified context file.
type ContextChange struct {
	Name    string
	ModTime time.Time
}

// CodeSummary summarizes code changes since the reference time.
type CodeSummary struct {
	CommitCount int
	LatestMsg   string
	Dirs        []string
	Authors     []string
}
