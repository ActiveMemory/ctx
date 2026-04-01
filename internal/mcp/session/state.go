//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import "time"

// NewState creates a new session state for the given context directory.
//
// Parameters:
//   - contextDir: Path to the project context directory
//
// Returns:
//   - *State: Initialized session state with empty counters
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
