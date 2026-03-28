//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

// Context-size event names.
const (
	// Suppressed is the event name for suppressed prompts.
	Suppressed = "suppressed"
	// Silent is the event name for silent (no-action) prompts.
	Silent = "silent"
	// Checkpoint is the event name for context checkpoint emissions.
	Checkpoint = "checkpoint"
	// WindowWarning is the event name for context window warning emissions.
	WindowWarning = "window-warning"
)
