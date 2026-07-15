//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import "time"

// PendingEntry records a context reference that has been staged for
// attachment to the next git commit.
//
// Fields:
//   - Ref: The staged context reference
//   - Timestamp: When the reference was staged
type PendingEntry struct {
	Ref       string    `json:"ref"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryEntry records the context references that were attached to a
// specific git commit.
//
// Fields:
//   - Commit: The commit hash the refs were attached to
//   - Refs: The context references attached to the commit
//   - Message: The commit subject line
//   - Timestamp: When the attachment was recorded
type HistoryEntry struct {
	Commit    string    `json:"commit"`
	Refs      []string  `json:"refs"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// OverrideEntry allows an explicit context association to be attached
// to a commit after the fact, replacing any automatically recorded refs.
//
// Fields:
//   - Commit: The commit hash to override refs for
//   - Refs: The replacement context references
//   - Timestamp: When the override was recorded
type OverrideEntry struct {
	Commit    string    `json:"commit"`
	Refs      []string  `json:"refs"`
	Timestamp time.Time `json:"timestamp"`
}

// ResolvedRef holds the result of resolving a raw context reference
// (e.g. "T-3", "D-1", "L-5") to its full details.
//
// Fields:
//   - Raw: The reference as written (e.g. "T-3")
//   - Type: Reference kind (task, decision, learning, ...)
//   - Number: Numeric ID parsed from the reference (0 if none)
//   - Title: Resolved entry title (empty if not found)
//   - Detail: Additional resolved detail (empty if not found)
//   - Found: True when the reference resolved to a known entry
type ResolvedRef struct {
	Raw    string
	Type   string
	Number int
	Title  string
	Detail string
	Found  bool
}
