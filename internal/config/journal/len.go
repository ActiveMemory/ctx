//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// Recall show/list display limits.
const (
	// PreviewMaxTurns is the maximum number of user turns shown in
	// the conversation preview of recall show.
	PreviewMaxTurns = 5
	// PreviewMaxTextLen is the maximum character length for a single
	// turn in the conversation preview.
	PreviewMaxTextLen = 100
	// SlugMaxLen is the maximum display length for session slugs in
	// recall list output.
	SlugMaxLen = 36
	// SessionIDShortLen is the prefix length for short session IDs
	// in summary output.
	SessionIDShortLen = 8
	// SessionIDHintLen is the prefix length for session IDs in
	// disambiguation hints (longer than short for uniqueness).
	SessionIDHintLen = 12
)
