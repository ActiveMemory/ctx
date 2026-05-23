//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"strings"

	"golang.org/x/text/unicode/norm"

	cfgI18n "github.com/ActiveMemory/ctx/internal/config/i18n"
)

// MatchKey returns a normalized comparison key for s:
// case-folded via [Fold], NFKD-decomposed, with general-
// purpose combining diacritical marks
// (cfgI18n.CombiningMarksLatinStart..End) stripped. The
// result is suitable for "user-intent" case- and
// diacritic-insensitive matching where the comparison
// should ignore casual keyboard variation (German
// `Straße` matches `STRASSE`; French `café` matches
// `cafe`; Turkish `İPTAL` matches `iptal`).
//
// Use MatchKey for placeholder/keyword vocabulary lookup,
// fuzzy search, and other user-input matching. Use [Fold]
// when you need Unicode-precise comparison that preserves
// linguistic distinctions (`İ` ≠ `i`, `ü` ≠ `u`).
//
// # What gets collapsed
//
//   - Latin diacritics: ü→u, é→e, ñ→n, ç→c, š→s, …
//   - Turkish dotted-I: İ→i (via Fold's i+U+0307 then
//     stripping U+0307).
//   - German sharp-s: ß→ss (via Fold).
//   - Greek final-sigma: ς→σ (via Fold with
//     HandleFinalSigma).
//   - Vietnamese horn: ơ→o, ư→u (U+031B is in the Latin
//     combining-marks block).
//
// # What stays distinct
//
//   - Arabic combining marks (U+0610–U+06ED): hamza,
//     shadda, fatha, kasra, damma, etc.
//   - Indic vowel signs (Bengali U+0980–U+09FF,
//     Devanagari U+0900–U+097F, Tamil, Telugu, …).
//   - Hebrew niqqud (U+0591–U+05C7).
//   - CJK ideographs (no combining marks; no false
//     collapse).
//
// # Limits
//
// MatchKey is a comparison primitive, not a tokenizer or
// normalizer. It does not strip whitespace, transliterate
// scripts (ä→ae is not done), or collapse digraphs that
// aren't already in NFKD (Polish ł stays distinct from l).
//
// Parameters:
//   - s: the input string. May be empty; may contain any
//     valid UTF-8.
//
// Returns:
//   - string: the normalized matching key.
func MatchKey(s string) string {
	decomposed := norm.NFKD.String(Fold(s))
	// Fast path: most inputs (pure ASCII, CJK, undecorated
	// non-Latin scripts) carry no runes in the Latin
	// combining-marks range, so we can return the
	// decomposed string directly without building.
	needsStrip := false
	for _, r := range decomposed {
		if r >= cfgI18n.CombiningMarksLatinStart &&
			r <= cfgI18n.CombiningMarksLatinEnd {
			needsStrip = true
			break
		}
	}
	if !needsStrip {
		return decomposed
	}
	var b strings.Builder
	b.Grow(len(decomposed))
	for _, r := range decomposed {
		if r >= cfgI18n.CombiningMarksLatinStart &&
			r <= cfgI18n.CombiningMarksLatinEnd {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
