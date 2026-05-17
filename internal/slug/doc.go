//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package slug generates **URL-safe, filesystem-safe
// identifiers** from human-readable input. Used by the
// journal importer to derive entry filenames and by the
// kb topic-page scaffolder to derive folder slugs.
//
// # Public Surface
//
//   - **[FromTitle](title)**: strict kebab-case slug.
//     Lowercases, replaces every non-alphanumeric run with
//     a single hyphen, trims, and truncates on a word
//     boundary at [journal.TitleSlugMaxLen]. Idempotent.
//   - **[CleanTitle](title)**: normalises a display title
//     for storage in YAML frontmatter (whitespace
//     collapsing + length cap). Pairs with FromTitle.
//   - **[ForTitle](session, existing)**: fallback chain
//     that picks the best title-derived slug for a
//     journal session (enriched title → first user msg →
//     Claude Code slug → short ID).
//   - **[Path](s)**: kebab-case slug that preserves `/` so
//     vendor-namespaced kb topic slugs survive
//     normalisation (e.g. `cursor/hooks`).
//
// # Stability Contract
//
// The slug for a given (title, dedup-context) pair is
// **deterministic**: re-running the importer against
// the same source produces the same slug. This is
// what makes the importer idempotent and what lets
// `git diff` show a meaningful patch when an entry
// is re-enriched.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
package slug
