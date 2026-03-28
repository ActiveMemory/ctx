//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import "time"

// State tracks per-context-dir advisory state.
//
// Session state is keyed by contextDir on the Server struct. It tracks
// tool call counts, entry additions, and pending context updates that
// need human review before persisting.
//
// Thread-safety: State is only accessed from the main request
// loop (single goroutine). If future work introduces concurrent access,
// a mutex should be added here.
type State struct {
	contextDir       string
	ToolCalls        int
	AddsPerformed    map[string]int
	sessionStartedAt time.Time
	PendingFlush     []PendingUpdate
}

// PendingUpdate represents a context update awaiting human confirmation.
type PendingUpdate struct {
	Type     string
	Content  string
	Attrs    map[string]string
	QueuedAt time.Time
}
