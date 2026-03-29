//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load_gate

const (
	// PrefixCtxLoaded is the filename prefix for session-loaded marker files.
	PrefixCtxLoaded = "ctx-loaded-"
	// EventContextLoadGate is the event name for context load gate hook events.
	EventContextLoadGate = "context-load-gate"
	// ContextLoadSeparatorChar is the character used for header/footer separators.
	ContextLoadSeparatorChar = "="
	// ContextLoadSeparatorWidth is the width of header/footer separator lines.
	ContextLoadSeparatorWidth = 80
	// ContextLoadIndexSuffix is the suffix appended to filenames
	// for index entries.
	ContextLoadIndexSuffix = " (idx)"
)
