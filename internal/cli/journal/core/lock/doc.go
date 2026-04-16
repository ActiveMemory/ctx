//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package lock manages **journal entry lock state** — the
// `locked: true` frontmatter flag that protects an enriched
// journal entry from being clobbered by a re-import of its
// raw source session.
//
// Locking is the journal pipeline's "do not touch" affordance.
// Without it, every `ctx journal import --regenerate` would
// risk overwriting the careful edits an author made to an
// enriched entry. With it, the importer sees `locked: true`,
// skips that file, and reports it in the import summary.
//
// # The Surface
//
//   - **[MatchJournalFiles](dir, pattern)** — finds journal
//     files matching a CLI pattern (slug, date, ID, glob).
//     Used by `ctx journal lock <pattern>` and `ctx journal
//     unlock <pattern>` to expand a pattern to a concrete
//     list of files. Pattern semantics match what the
//     user-facing CLI documents.
//   - **[MultipartBase](filename)** — extracts the base
//     name from a multipart filename (e.g.
//     `2026-04-12-foo--part2.md` → `2026-04-12-foo`). The
//     lock state for a multipart entry lives on the **base
//     part**, and other parts inherit it.
//   - **[UpdateFrontmatter](path, lock)** — atomic update
//     of the `locked:` field in a file's YAML frontmatter.
//     Adds the field if missing; removes it when `lock` is
//     false (rather than writing `locked: false`, which
//     would still bypass the importer's omit-default
//     check).
//
// # State File Sync
//
// The lock state can also be read from
// `.context/journal/.state.json` (per
// [internal/journal/state]). Frontmatter is the source of
// truth; the state file is a denormalized index for fast
// queries from `ctx journal sync` and the importer. The
// `ctx journal sync` command (in
// [internal/cli/journal/cmd/sync]) reconciles drift in
// either direction.
//
// # Concurrency
//
// All operations are file-local and hold the file open
// only for the duration of the read+write. Concurrent
// invocations against different files never race;
// concurrent updates to the same file would race on the
// final write (no per-file locking is implemented — the
// CLI is single-process anyway).
package lock
