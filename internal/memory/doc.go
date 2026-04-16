//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory bridges Claude Code's per-project
// auto memory file (MEMORY.md) into a project's
// .context/ directory so that memory written by
// Claude becomes git-tracked, version-controlled,
// drift-checkable, and importable into the structured
// context files (DECISIONS / LEARNINGS / CONVENTIONS).
//
// # The Problem It Solves
//
// Claude Code's auto memory lives at
// ~/.claude/projects/<slug>/memory/MEMORY.md. That
// path is not in the project repo (no peer review,
// no git history), is per-machine (teammates see
// different files), and silently grows as Claude
// takes notes, drifting from .context/ over time.
//
// # Pipeline Stages
//
//   - **discover** ([DiscoverPath], [ProjectSlug]):
//     encodes the project root into Claude Code's
//     slug format and resolves the auto-memory file.
//   - **sync** ([Sync], [Archive]): copies the
//     source into .context/memory/mirror.md, archiving
//     the previous mirror before overwrite.
//   - **diff** ([Diff], [HasDrift]): line-level diff
//     between source and mirror; surfaces what Claude
//     wrote since the last sync.
//   - **parse** ([Entries]): splits MEMORY.md
//     content into discrete [Entry] blocks by headers,
//     blank lines, and list items.
//   - **classify** ([Classify]): routes each entry
//     to the matching .context/ file based on keyword
//     heuristics from .ctxrc classify_rules.
//   - **promote** ([Promote]): writes a classified
//     entry to its target .context/ file.
//   - **publish** ([Publish], [SelectContent],
//     [MergePublished], [RemovePublished]): the
//     inverse direction: promotes .context/ entries
//     into MEMORY.md so future Claude sessions see
//     them up front.
//
// # State Tracking
//
// Sync and import state lives in
// .context/state/memory-import.json ([LoadState],
// [SaveState]). The [State] struct tracks last-synced
// timestamps, imported entry hashes ([EntryHash]),
// and import/publish progress.
//
// # Concurrency and Idempotency
//
// All operations are read or write-once with no
// long-lived goroutines. [Sync] is idempotent: an
// unchanged source produces no archive entry and no
// mirror write.
package memory
