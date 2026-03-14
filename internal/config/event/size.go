//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

// Context-size event names.
const (
	// EventSuppressed is the event name for suppressed prompts.
	EventSuppressed = "suppressed"
	// EventSilent is the event name for silent (no-action) prompts.
	EventSilent = "silent"
	// EventCheckpoint is the event name for context checkpoint emissions.
	EventCheckpoint = "checkpoint"
	// EventWindowWarning is the event name for context window warning emissions.
	EventWindowWarning = "window-warning"
)
