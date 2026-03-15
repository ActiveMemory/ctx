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

// NewState creates a new session state for the given context directory.
func NewState(contextDir string) *State {
	return &State{
		contextDir:       contextDir,
		AddsPerformed:    make(map[string]int),
		sessionStartedAt: time.Now(),
	}
}

// RecordToolCall increments the tool call counter.
func (ss *State) RecordToolCall() {
	ss.ToolCalls++
}

// RecordAdd increments the add counter for the given entry type.
func (ss *State) RecordAdd(entryType string) {
	ss.AddsPerformed[entryType]++
}

// QueuePendingUpdate adds an update to the pending flush queue.
func (ss *State) QueuePendingUpdate(update PendingUpdate) {
	ss.PendingFlush = append(ss.PendingFlush, update)
}

// PendingCount returns the number of pending updates.
func (ss *State) PendingCount() int {
	return len(ss.PendingFlush)
}
