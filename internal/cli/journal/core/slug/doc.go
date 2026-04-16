//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package slug generates **URL-safe, filesystem-safe
// identifiers** from session titles and handles the
// deduplication logic that keeps two sessions with the
// same title from colliding on disk.
//
// Slugs are how journal entries are addressed throughout
// ctx: filenames are `YYYY-MM-DD-<slug>.md`, links use
// the slug as the path component, and `ctx journal source
// --show <slug>` looks up by slug.
//
// # Public Surface
//
//   - **[FromTitle](title)** — converts a title to a
//     lowercase, hyphenated, alphanumeric slug.
//     Strips punctuation, collapses runs of separators,
//     trims leading/trailing hyphens. Idempotent.
//   - **[CleanTitle](title)** — strips non-alphanumeric
//     characters from a display title (kept for the
//     filename's human-readable suffix when one is
//     wanted in addition to the slug).
//   - **[ForTitle](title, existing)** — the dedup-aware
//     wrapper: produces both a slug and a cleaned
//     title, appending `-2`, `-3`, etc. when the
//     base slug already exists in `existing`. Used by
//     the importer when two sessions share a topic
//     summary.
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
