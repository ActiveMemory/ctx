//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package i18n holds the configuration constants used by
// `internal/i18n` — Unicode block boundaries and other
// numeric constants that would trip the magic-value AST
// audit if they lived inline in the implementation.
//
// # Public Surface
//
//   - [CombiningMarksLatinStart]: U+0300, start of the
//     Latin combining-diacritics block.
//   - [CombiningMarksLatinEnd]: U+036F, end of the same
//     block.
//
// Both bound the range that `internal/i18n.MatchKey`
// strips to produce a diacritic-insensitive comparison
// key.
package i18n
