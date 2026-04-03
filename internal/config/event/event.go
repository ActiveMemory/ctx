//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

// Session lifecycle event types.
const (
	// TypeStart is the session start event.
	TypeStart = "start"
	// TypeEnd is the session end event.
	TypeEnd = "end"
)

// Notify event type constants.
const (
	// TypeTest is the event type for test notifications.
	TypeTest = "test"
	// TestMessage is the payload message for test notifications.
	TestMessage = "Test notification from ctx"
)

// Event categories for log grouping.
const (
	// CategorySession groups session lifecycle events.
	CategorySession = "session"
)

// Template variable keys for session events.
const (
	// VarCaller is the template variable for the calling editor.
	VarCaller = "Caller"
)

// Events display configuration.
const (
	// MessageMaxLen is the maximum character length for event messages
	// in human-readable output before truncation.
	MessageMaxLen = 60
	// HookFallback is the placeholder displayed when no hook name
	// can be determined from an event payload.
	HookFallback = "-"
	// TruncationSuffix is appended to truncated event messages.
	TruncationSuffix = "..."
)
