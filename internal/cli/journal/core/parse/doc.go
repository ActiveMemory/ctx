//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parse scans a journal directory and turns each
// markdown file into a typed [entity.JournalEntry] by
// reading and validating its YAML frontmatter — the
// upstream of every site-builder, Obsidian-exporter, MOC
// generator, and lock-state syncer in the journal pipeline.
//
// # Public Surface
//
//   - **[ScanJournalEntries](dir)** — walks `dir`,
//     parses every `*.md` file, and returns
//     `[]*JournalEntry` plus an error slice for files
//     that failed to parse. The walk continues past
//     bad files so a single malformed entry does not
//     abort the whole scan.
//   - **[JournalEntry](path)** — parses one file by
//     path. Used by single-entry callers (the
//     `--show <slug>` lookup, the lock CLI, and the
//     drift checker).
//
// # Frontmatter Schema
//
// Each entry's frontmatter must satisfy the journal
// schema documented in
// `internal/entity/journal.go.JournalFrontmatter`:
// `id`, `date`, `title`, `slug`, optional `topics`,
// optional `locked`, optional `enriched`, optional
// `part` / `parts` for multipart entries. Unknown
// fields are preserved (round-trip safe).
//
// # Concurrency
//
// Stateless and filesystem-bound. Concurrent calls
// against the same directory each pay the full read
// cost.
//
// # Related Packages
//
//   - [internal/cli/journal/cmd/site] /
//     [internal/cli/journal/cmd/obsidian]   — chief
//     consumers of [ScanJournalEntries].
//   - [internal/cli/journal/core/lock]       — uses
//     [JournalEntry] for single-file lock checks.
//   - [internal/cli/system/cmd/check_journal] — uses
//     the count to nudge about pending imports.
//   - [internal/parse]                       — supplies
//     [SplitFrontmatter] used internally.
package parse
