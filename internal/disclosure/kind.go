//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
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

// ThemeDir returns the context-relative subdirectory that holds a kind's
// theme files (<noun>/<slug>.md). It is false only for an unknown kind,
// so the mover refuses rather than write entry bodies to a guessed path.
//
// Parameters:
//   - k: the root kind
//
// Returns:
//   - string: the theme-file subdirectory (meaningful only when ok)
//   - bool: true for every kind the mover digests
func ThemeDir(k Kind) (string, bool) {
	switch k {
	case KindLearning:
		return cfgDisc.ThemeDirLearning, true
	case KindDecision:
		return cfgDisc.ThemeDirDecision, true
	case KindConvention:
		return cfgDisc.ThemeDirConvention, true
	default:
		return "", false
	}
}

// EntryPrefix returns the line prefix that opens an entry heading in a
// kind's staging zone: "## [" for the timestamped kinds (learning,
// decision), "## " for the curated prose sections of a convention root.
//
// Parametrizing this prefix is what collapses the two former parse paths
// into one — the structural difference between the kinds is the prefix
// and the identity (timestamp+title vs title alone), not the layout.
//
// Callers scanning a convention root must skip [cfgDisc.HeadingThemes],
// which shares the "## " prefix but is structure, not an entry.
//
// Parameters:
//   - k: the root kind
//
// Returns:
//   - string: the kind's entry-heading line prefix
func EntryPrefix(k Kind) string {
	switch k {
	case KindLearning, KindDecision:
		return cfgDisc.EntryLinePrefix
	case KindConvention:
		return cfgDisc.SectionLinePrefix
	default:
		return cfgDisc.EntryLinePrefix
	}
}
