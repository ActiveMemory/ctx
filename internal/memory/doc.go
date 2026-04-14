//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory bridges Claude Code's per-project **auto memory**
// file (`MEMORY.md`) into a project's `.context/` directory so that
// memory written by Claude — outside the repo, under
// `~/.claude/projects/<slug>/memory/` — becomes git-tracked,
// version-controlled, drift-checkable, and importable into the
// structured context files (DECISIONS / LEARNINGS / CONVENTIONS).
//
// # The Problem It Solves
//
// Claude Code's auto memory lives at
// `~/.claude/projects/<slug>/memory/MEMORY.md`. That path:
//
//   - is **not** in the project repo (no peer review, no git
//     history, no diff-on-PR);
//   - is **per-machine** (a teammate working on the same project
//     sees a different file);
//   - **silently grows** as Claude takes notes during sessions,
//     drifting from `.context/` content over time.
//
// This package makes that file a first-class citizen.
//
// # Pipeline Stages
//
// Stages are wired to the `ctx memory` CLI surface:
//
//   - **discover** ([Discover], [discover.go]) — encodes the
//     project root path into Claude Code's slug format
//     (absolute path, `/` → `-`, `-` prefix) and resolves the
//     auto-memory file. Returns an actionable error if the
//     auto-memory file does not exist.
//   - **status** — reports whether the source file exists,
//     when the last sync happened, and whether drift is
//     present (size or content).
//   - **sync** ([Mirror], [mirror.go]) — copies the source
//     into `.context/memory/mirror.md`. Previous mirror is
//     archived under `.context/memory/archive/<ts>.md` before
//     overwrite, so history is preserved.
//   - **diff** ([diff.go]) — line-level diff between source
//     and mirror; surfaces what Claude wrote since the last
//     sync.
//   - **import** ([extract.go], [classify.go]) — parses the
//     mirror into discrete entries and routes each one to the
//     matching `.context/` file based on keyword heuristics
//     ([Classify]); rules come from `.ctxrc.classify_rules`
//     with built-in fallbacks. The user gets a preview before
//     anything is written.
//   - **publish** ([publish.go], [promote.go]) — the inverse
//     direction: promotes a `.context/` entry into the
//     auto-memory file so future Claude sessions get it
//     up front.
//
// # State Tracking
//
// Sync and import state lives in
// `.context/state/memory-import.json` (see [state.go]):
// last-synced timestamps, last-imported entry hashes, and the
// drift signal computed from comparing source-vs-mirror file
// sizes. The drift hook
// (`internal/cli/system/cmd/check_memory_drift`) reads this
// state and nudges the user when `MEMORY.md` has changed since
// last sync.
//
// # Concurrency and Idempotency
//
// All operations are **read or write-once** — no long-lived
// goroutines. [Mirror] is idempotent: an unchanged source
// produces no archive entry and no mirror write. [Discover]
// caches its result in process memory but the cache is
// keyed on `projectRoot`, so different projects do not
// collide.
//
// # Related Packages
//
//   - [internal/cli/memory]   — the `ctx memory` CLI surface
//     (status, sync, diff, import, publish, unpublish).
//   - [internal/config/memory] — slug-format, path constants,
//     classification-rule schema.
//   - [internal/drift]        — consumes the drift state to
//     produce the user-facing nudge.
//   - [internal/cli/system/cmd/check_memory_drift] — the
//     hook that fires the nudge.
//   - [internal/write/memory] — terminal output for the CLI
//     subcommands.
package memory
