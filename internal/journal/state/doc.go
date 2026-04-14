//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package state manages the **journal processing state**
// stored in `.context/journal/.state.json` — a denormalized
// JSON index that tracks, for each raw session source, where
// it is in the import → normalize → enrich → wrap pipeline
// and whether it is locked against re-import.
//
// # Why an External State File
//
// The original design embedded markers in the journal files
// themselves (`<!-- normalized: ... -->`, etc.). That broke
// when a journal entry's body legitimately contained one of
// those marker strings — the parser saw a false positive and
// concluded the entry had been processed when it had not.
//
// Moving state out of the file body fixes the false-positive
// problem and gives the importer a fast index it can scan
// without parsing every entry.
//
// # The State File Shape
//
// `.state.json` is a `map[sourceID]Record` where each
// [Record] tracks:
//
//   - **stage** — current pipeline stage (one of
//     [ValidStages]: imported, normalized, enriched,
//     wrapped, indexed).
//   - **locked** — true when the entry is protected
//     from re-import regeneration.
//   - **part** — for multipart entries, which part this
//     record refers to.
//   - **filename** — the on-disk filename the importer
//     produced for this source.
//
// # Public Surface
//
//   - **[Load](journalDir)** — reads `.state.json` and
//     returns the deserialized map. Returns an empty
//     map (not an error) when the file is missing —
//     fresh projects have no state yet.
//   - **[ValidStages]** — canonical stage names in
//     pipeline order. Stage advancement is forward-only
//     in normal flow; re-import (`--regenerate`) resets
//     to "imported".
//
// # Sync With Frontmatter
//
// Frontmatter is the source of truth for `locked:`; the
// state file is a denormalized cache. `ctx journal sync`
// reconciles drift between the two so users who edit
// frontmatter directly see the importer respect the
// change on next run.
//
// # Concurrency
//
// File reads are scoped per call. Writes from the
// importer use the atomic-rename pattern so a partial
// write never produces a malformed JSON file.
//
// # Related Packages
//
//   - [internal/cli/journal/cmd/importer] — the chief
//     reader/writer.
//   - [internal/cli/journal/core/lock]    — the locking
//     helpers; mutate frontmatter and re-sync state.
//   - [internal/cli/journal/cmd/sync]     — explicit
//     reconciliation between frontmatter and state.
//   - [internal/cli/system/cmd/check_journal] — reads
//     state to count unimported sources.
package state
