//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task implements the **`ctx task`** command group:
// task completion, archival, and snapshots, the lifecycle
// operations on `TASKS.md`.
//
// `TASKS.md` is the project's living checklist. Phase
// headers are constitutional structure (never moved or
// renamed); items are append-only with status flips
// (`[ ]` → `[x]` or `[-]`) when work completes or is
// skipped. This package owns the safe transitions.
//
// # Subcommands
//
//   - **`ctx task complete [number|text]`**: flips a
//     task from `[ ]` to `[x]`. Match by phase-relative
//     number (e.g. `3`), partial text, or full text in
//     quotes. See [internal/cli/task/cmd/complete].
//   - **`ctx task archive [--dry-run]`**: moves
//     completed top-level tasks into a dated
//     `.context/archive/tasks-YYYY-MM-DD.md` file,
//     preserving phase structure. See
//     [internal/cli/task/cmd/archive] (delegates to
//     [internal/tidy.WriteArchive]).
//   - **`ctx task snapshot [name]`**: copies the
//     current TASKS.md verbatim to
//     `.context/archive/snapshots/<ts>-<name>.md`. No
//     mutation of the source. Used before a major
//     restructure to give the user a known-good
//     restore point.
//
// # Constitutional Invariants
//
// The CONSTITUTION.md rules:
//
//   - **Tasks stay in their Phase section permanently**.
//   - **Phase headers are never removed or renamed**.
//   - **Tasks are never deleted**, only marked
//     `[x]` (completed) or `[-]` (skipped).
//   - **Archival ≠ deletion**: archived tasks land in
//     the archive file, not `/dev/null`.
//
// This package enforces all four: `ctx task complete`
// uses status flips, never moves. `ctx task archive`
// uses [internal/tidy] which preserves phase structure
// in the archive output.
//
// # Concurrency
//
// Filesystem-bound and stateless. Single-process
// assumption.
package task
