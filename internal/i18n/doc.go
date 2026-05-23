//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package i18n provides Unicode-correct primitives for
// internationalization-sensitive string operations.
// Exports two case-insensitive comparison primitives with
// different contracts: [Fold] is the strict Unicode
// case-fold; [MatchKey] is the ergonomic, diacritic-
// insensitive variant for casual user-input matching.
//
// # Why this package exists
//
// `strings.ToLower` is locale-naive and produces output
// unsuitable for cross-locale comparison. Turkish input
// (İ→i̇), German input (ß→ss), Greek input (final-sigma),
// and many others fold incorrectly. Code that searches,
// classifies, or matches user-supplied strings needs
// Unicode-correct primitives.
//
// # Public Surface
//
//   - **[Fold](s)**: returns the Unicode case-folded form
//     of s. Backed by `golang.org/x/text/cases.Fold` with
//     `HandleFinalSigma(true)`. Byte-identical to
//     strings.ToLower for ASCII input. Preserves
//     linguistic distinctions: `İ` ≠ `i`, `ü` ≠ `u`.
//     Use when you want Unicode-precise comparison.
//   - **[MatchKey](s)**: returns a casual-comparison key
//     for s: Fold + NFKD + strip Latin combining marks
//     (U+0300–U+036F). Collapses Turkish dotted-I,
//     German umlauts, French accents, Vietnamese horns,
//     and similar Latin/general diacritics. Preserves
//     script-essential marks for Arabic (hamza, niqqud),
//     Indic (vowel signs), Hebrew (niqqud), and CJK
//     (no diacritics). Use for placeholder/keyword
//     vocabulary lookup and other user-intent matching.
//
// # Picking the right primitive
//
// Rule of thumb: if your matcher compares user input
// against a vocabulary list and the user might type with
// or without diacritics, use MatchKey. If you need to
// preserve Unicode-defined linguistic distinctions
// (parsing, deduplication of normalized identifiers,
// security-relevant comparison), use Fold.
//
// # Enforcement
//
// The `internal/compliance` package contains an AST test
// (`TestNoDirectStringsToLower`) that fails on any direct
// `strings.ToLower` call in the codebase outside this
// package itself. No allowlist, no per-package opt-out.
// New code that wants case-insensitive comparison must
// call [Fold] or [MatchKey].
package i18n
