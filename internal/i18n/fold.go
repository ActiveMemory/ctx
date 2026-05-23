//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"golang.org/x/text/cases"
)

// foldCaser is the package-singleton Unicode case-folder.
// `cases.Caser` is safe for concurrent use per the upstream
// docs, so a single instance is reused across calls.
var foldCaser = cases.Fold(cases.HandleFinalSigma(true))

// Fold returns the Unicode case-folded form of s, suitable
// for case-insensitive string comparison. This is the
// project-mandated replacement for strings.ToLower in any
// comparison context — `internal/compliance` bans direct
// strings.ToLower calls outside this package.
//
// Why fold and not lower? Lowercasing is locale-dependent
// (Turkish dotted/dotless I, etc.) and produces strings
// unsuitable for cross-locale comparison. Unicode case
// folding is the stdlib-recommended primitive for
// case-insensitive matching:
//
//   - Turkish İ → i̇ + COMBINING DOT ABOVE, then folded
//     such that İ and i compare equal regardless of
//     locale.
//   - German ß → ss, so "STRASSE" and "Straße" compare
//     equal.
//   - Greek ς (final sigma) and σ fold to the same form
//     when HandleFinalSigma is on (it is, here).
//
// For pure ASCII input, Fold is byte-identical to
// strings.ToLower, so swapping is behavior-preserving on
// the ASCII paths that the codebase has historically
// relied on.
//
// Parameters:
//   - s: the input string. May be empty; may contain any
//     valid UTF-8.
//
// Returns:
//   - string: the case-folded form. Length may differ from
//     len(s) when folding expands characters (e.g. ß→ss).
func Fold(s string) string {
	return foldCaser.String(s)
}
