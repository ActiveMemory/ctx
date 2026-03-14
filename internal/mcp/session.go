//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import "time"

// sessionState tracks per-context-dir advisory state.
//
// Session state is keyed by contextDir on the Server struct. It tracks
// tool call counts, entry additions, and pending context updates that
// need human review before persisting.
type sessionState struct {
	contextDir       string
	toolCalls        int
	addsPerformed    map[string]int
	sessionStartedAt time.Time
	pendingFlush     []PendingUpdate
}

// PendingUpdate represents a context update awaiting human confirmation.
type PendingUpdate struct {
	Type     string
	Content  string
	Attrs    map[string]string
	QueuedAt time.Time
}

// newSessionState creates a new session state for the given context directory.
func newSessionState(contextDir string) *sessionState {
	return &sessionState{
		contextDir:       contextDir,
		addsPerformed:    make(map[string]int),
		sessionStartedAt: time.Now(),
	}
}

// recordToolCall increments the tool call counter.
func (ss *sessionState) recordToolCall() {
	ss.toolCalls++
}

// recordAdd increments the add counter for the given entry type.
func (ss *sessionState) recordAdd(entryType string) {
	ss.addsPerformed[entryType]++
}

// queuePendingUpdate adds an update to the pending flush queue.
func (ss *sessionState) queuePendingUpdate(update PendingUpdate) {
	ss.pendingFlush = append(ss.pendingFlush, update)
}

// pendingCount returns the number of pending updates.
func (ss *sessionState) pendingCount() int {
	return len(ss.pendingFlush)
}
