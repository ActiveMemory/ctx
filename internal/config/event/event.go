//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

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
