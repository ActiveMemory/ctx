//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package disclosure holds the structural vocabulary for
// progressive-disclosure roots (see specs/progressive-disclosure.md).
//
// # Domain
//
// A bounded root is delimited by ATX headings whose exact text is
// load-bearing: [HeadingThemes] ("## Themes") marks the themes region
// and the lower bound of the staging zone, for every kind.
//
// Within the staging zone, one line prefix opens each entry, and it is
// the only structural difference between the kinds: [EntryLinePrefix]
// ("## [") for the timestamped LEARNINGS/DECISIONS entries, and
// [SectionLinePrefix] ("## ") for the curated prose sections of a
// CONVENTIONS root. [ThemeDirLearning], [ThemeDirDecision], and
// [ThemeDirConvention] name the per-kind subdirectories of the context
// directory that hold theme files.
//
// These constants are a single source of truth: the parser
// (internal/disclosure), the validate precondition, and the add-path
// layout proofs all key on these exact strings.
//
// # Concurrency
//
// Compile-time constants; safe for concurrent use.
package disclosure
