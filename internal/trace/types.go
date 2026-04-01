//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import "time"

// PendingEntry records a context reference that has been staged for
// attachment to the next git commit.
type PendingEntry struct {
	Ref       string    `json:"ref"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryEntry records the context references that were attached to a
// specific git commit.
type HistoryEntry struct {
	Commit    string    `json:"commit"`
	Refs      []string  `json:"refs"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// OverrideEntry allows an explicit context association to be attached
// to a commit after the fact, replacing any automatically recorded refs.
type OverrideEntry struct {
	Commit    string    `json:"commit"`
	Refs      []string  `json:"refs"`
	Timestamp time.Time `json:"timestamp"`
}

// ResolvedRef holds the result of resolving a raw context reference
// (e.g. "T-3", "D-1", "L-5") to its full details.
type ResolvedRef struct {
	Raw    string
	Type   string
	Number int
	Title  string
	Detail string
	Found  bool
}
