//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

import "time"

// MCPSession tracks advisory state for one MCP server run.
//
// The MCP server keeps a single MCPSession for the lifetime of a
// context directory. It records tool call counts, entry additions,
// and pending context updates that need human review before they
// persist. Governance helpers in the handler package read this
// state to decide which advisory warnings to append to tool
// responses.
//
// Thread-safety: MCPSession is only touched from the main request
// loop (single goroutine). If future work introduces concurrent
// access, a mutex should be added here.
//
// Fields:
//   - ToolCalls: Total tool invocations in this session
//   - AddsPerformed: Entry additions by type (decision, learning, etc.)
//   - SessionStartedAt: Session start timestamp
//   - PendingFlush: Updates awaiting human confirmation
//   - SessionStarted: Whether session_event:start has fired
//   - ContextLoaded: Whether context files have been read
//   - LastDriftCheck: Timestamp of most recent drift check
//   - LastContextWrite: Timestamp of most recent .context write
//   - CallsSinceWrite: Tool calls since last .context write
type MCPSession struct {
	ToolCalls        int
	AddsPerformed    map[string]int
	SessionStartedAt time.Time
	PendingFlush     []PendingUpdate

	SessionStarted   bool
	ContextLoaded    bool
	LastDriftCheck   time.Time
	LastContextWrite time.Time
	CallsSinceWrite  int
}

// PendingUpdate represents a context update awaiting human confirmation.
//
// Fields:
//   - Type: Update type (decision, learning, task, convention)
//   - Content: Entry text
//   - Attrs: Optional attributes (context, rationale, etc.)
//   - QueuedAt: When this update was queued
type PendingUpdate struct {
	Type     string
	Content  string
	Attrs    map[string]string
	QueuedAt time.Time
}

// NewMCPSession creates a fresh MCP session with empty counters.
//
// Returns:
//   - *MCPSession: initialized session ready to record events
func NewMCPSession() *MCPSession {
	return &MCPSession{
		AddsPerformed:    make(map[string]int),
		SessionStartedAt: time.Now(),
	}
}

// RecordToolCall increments the tool call counter.
func (ss *MCPSession) RecordToolCall() {
	ss.ToolCalls++
}

// RecordAdd increments the add counter for the given entry type.
//
// Parameters:
//   - entryType: Context entry type (task, decision, etc.)
func (ss *MCPSession) RecordAdd(entryType string) {
	ss.AddsPerformed[entryType]++
}

// QueuePendingUpdate adds an update to the pending flush queue.
//
// Parameters:
//   - update: Update to enqueue for the next flush cycle
func (ss *MCPSession) QueuePendingUpdate(update PendingUpdate) {
	ss.PendingFlush = append(ss.PendingFlush, update)
}

// PendingCount returns the number of pending updates.
//
// Returns:
//   - int: Number of updates awaiting flush
func (ss *MCPSession) PendingCount() int {
	return len(ss.PendingFlush)
}

// RecordSessionStart marks the session as explicitly started and
// resets the session start timestamp.
//
// Called by the session_event tool when the agent reports a "start"
// event. Sets SessionStarted to true and captures the current wall
// time so governance checks can measure elapsed time.
func (ss *MCPSession) RecordSessionStart() {
	ss.SessionStarted = true
	ss.SessionStartedAt = time.Now()
}

// RecordContextLoaded marks context as loaded for this session.
//
// Called after the agent successfully loads context files (TASKS.md,
// DECISIONS.md, etc.). Suppresses the "context not loaded" governance
// warning that would otherwise appear on every tool response.
func (ss *MCPSession) RecordContextLoaded() {
	ss.ContextLoaded = true
}

// RecordDriftCheck records that a drift check was performed.
//
// Called after the agent runs ctx_drift. Updates the last-drift-check
// timestamp so governance helpers can determine whether a follow-up
// drift check is overdue.
func (ss *MCPSession) RecordDriftCheck() {
	ss.LastDriftCheck = time.Now()
}

// RecordContextWrite records that a .context/ write occurred.
//
// Called after successful ctx_add, ctx_complete, ctx_watch_update,
// or ctx_compact invocations. Captures the current wall time and
// resets the calls-since-write counter to zero.
func (ss *MCPSession) RecordContextWrite() {
	ss.LastContextWrite = time.Now()
	ss.CallsSinceWrite = 0
}

// IncrementCallsSinceWrite bumps the counter used for persist nudges.
//
// Called by the MCP server after every tool dispatch regardless of
// tool type. When the counter reaches the persist-nudge threshold,
// governance helpers begin emitting persist nudge warnings.
func (ss *MCPSession) IncrementCallsSinceWrite() {
	ss.CallsSinceWrite++
}
