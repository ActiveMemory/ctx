//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package i18n

// CombiningMarksLatinStart is the first code point of the
// Unicode "Combining Diacritical Marks" block. This block
// holds Latin and general-purpose diacritics: acute,
// grave, diaeresis, tilde, cedilla, caron, the Turkish
// combining dot above (U+0307), the Vietnamese combining
// horn (U+031B), and similar marks that ride on top of
// Latin-script bases. Script-essential marks for Arabic,
// Indic, Hebrew, and other scripts live in their own
// blocks and are *not* part of this range.
const CombiningMarksLatinStart = 0x0300

// CombiningMarksLatinEnd is the last code point of the
// Unicode "Combining Diacritical Marks" block. The
// [CombiningMarksLatinStart, CombiningMarksLatinEnd] range
// is what `internal/i18n.MatchKey` strips when producing
// a diacritic-insensitive comparison key — that's how
// `İPTAL` collapses to `iptal` and `Straße` collapses to
// `strasse` while Arabic / Indic / Hebrew script-essential
// marks stay distinct.
const CombiningMarksLatinEnd = 0x036F
