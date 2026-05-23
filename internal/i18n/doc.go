//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package i18n provides Unicode-correct primitives for
// internationalization-sensitive string operations.
// Today's sole export is [Fold], the project-mandated
// replacement for `strings.ToLower` in case-insensitive
// comparison contexts.
//
// # Why this package exists
//
// `strings.ToLower` is locale-naive and produces output
// unsuitable for cross-locale comparison. Turkish input
// (İ→i̇), German input (ß→ss), Greek input (final-sigma),
// and many others fold incorrectly. Code that searches,
// classifies, or matches user-supplied strings needs
// Unicode-correct folding.
//
// # Public Surface
//
//   - **[Fold](s)**: returns the Unicode case-folded form
//     of s. Backed by `golang.org/x/text/cases.Fold` with
//     `HandleFinalSigma(true)`. Byte-identical to
//     strings.ToLower for ASCII input.
//
// # Enforcement
//
// The `internal/compliance` package contains an AST test
// (`TestNoDirectStringsToLower`) that fails on any direct
// `strings.ToLower` call in the codebase outside this
// package itself. No allowlist, no per-package opt-out.
// New code that wants case-insensitive comparison must
// call [Fold].
package i18n
