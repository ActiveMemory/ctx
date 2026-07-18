//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgEntry "github.com/ActiveMemory/ctx/internal/config/entry"
)

// KindFor maps a canonical knowledge-file basename to its Kind. It is
// how the CLI decides which root it was handed; a non-knowledge file
// returns false so the caller can refuse rather than guess.
//
// Parameters:
//   - basename: a file's base name (e.g. "LEARNINGS.md")
//
// Returns:
//   - Kind: the matched kind (meaningful only when ok is true)
//   - bool: true when basename is a canonical knowledge file
func KindFor(basename string) (Kind, bool) {
	switch basename {
	case cfgCtx.Learning:
		return KindLearning, true
	case cfgCtx.Decision:
		return KindDecision, true
	case cfgCtx.Convention:
		return KindConvention, true
	default:
		return KindLearning, false
	}
}

// String returns the kind's name, matching the entry-type vocabulary
// ("learning" | "decision" | "convention"). Used for the Inspection's
// stable string Kind field.
//
// Returns:
//   - string: the kind name, or "" for an unknown kind
func (k Kind) String() string {
	switch k {
	case KindLearning:
		return cfgEntry.Learning
	case KindDecision:
		return cfgEntry.Decision
	case KindConvention:
		return cfgEntry.Convention
	default:
		return ""
	}
}
