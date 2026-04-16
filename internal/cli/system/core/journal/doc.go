//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal holds the **shared helpers** that the
// `check_journal` hook calls when deciding whether to nudge
// the user about unimported or unenriched session entries.
//
// The package is the *measurement layer*; the hook decides
// when the numbers warrant a nudge.
//
// # Public Surface
//
//   - **[NewestMtime](dir)**: returns the modification
//     time of the most recently changed file in `dir`,
//     or zero time when the directory is empty/missing.
//     Used to compare the journal directory's freshness
//     to the raw-source directory.
//   - **[CountNewerFiles](dir, since)**: counts files
//     in `dir` modified strictly after `since`. The
//     hook calls this with the last-known-import
//     timestamp to surface "N new sessions to import".
//   - **[CountUnenriched](journalDir)**: counts entries
//     in `journalDir` that have not been enriched
//     (frontmatter `enriched: false` or missing).
//     Surfaces the "N entries waiting for enrichment"
//     nudge.
//   - **[CheckStage](path, stage)**: predicate; is
//     `path`'s frontmatter at or past `stage`?
//     [internal/journal/state] supplies the canonical
//     stage strings.
//   - **[MarkStage](path, stage)**: atomically updates
//     the frontmatter `stage:` field. Used by the
//     enrichment pipeline to advance an entry through
//     normalize → enrich → wrap.
//
// # Concurrency
//
// All functions are filesystem-bound and stateless.
// [MarkStage] uses the same atomic-rename pattern as
// the rest of the journal pipeline so a partial write
// never leaves an entry in an undefined state.
package journal
