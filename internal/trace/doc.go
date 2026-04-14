//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace implements **commit context tracing** — the layer
// that links a git commit back to the decisions, learnings,
// conventions, tasks, and AI sessions that motivated it.
//
// The point is to make `git log -p` answer not just "what
// changed" but "*why* it changed", without forcing the developer
// to write hand-curated provenance every time. The package
// gathers context references from three sources at commit time,
// renders them as a structured git trailer, and persists a
// per-commit history record so the link survives even when the
// commit message is later squashed or rewritten.
//
// # Reference Format
//
// A "ref" is a short, parseable string that points at one
// concrete piece of context:
//
//   - `decision:12`   — DECISIONS.md entry #12
//   - `learning:7`    — LEARNINGS.md entry #7
//   - `convention:3`  — CONVENTIONS.md entry #3
//   - `task:8`        — TASKS.md item #8
//   - `session:abc`   — AI session ID `abc`
//   - `"free note"`   — quoted free-form note
//
// [parseRef] turns a string into (type, number, text); [Resolve]
// looks up the entry and returns a [ResolvedRef] populated with
// the entry title and a one-line detail preview.
//
// # The Three-Source Collection
//
// [Collect] runs at commit time (typically from a `prepare-
// commit-msg` hook) and gathers refs from three independent
// sources, **in this order**, then deduplicates while preserving
// first-occurrence order:
//
//  1. **Pending records** — refs that were explicitly staged
//     ahead of time via `ctx trace tag` and stored as
//     [PendingEntry] in `state/trace-pending.jsonl`. Cleared
//     after the commit lands.
//  2. **Staged file diffs** — [StagedRefs] runs `git diff
//     --cached` on each of DECISIONS.md, LEARNINGS.md,
//     CONVENTIONS.md and parses **added entries** into refs of
//     the matching type. For TASKS.md it parses **completed
//     tasks** (lines that flipped from `[ ]` to `[x]`). This is
//     the source that catches "I just wrote a new decision and
//     committed it" without any tagging.
//  3. **Working state** — [WorkingRefs] adds in-progress task
//     refs (from TASKS.md) plus an `session:<id>` ref derived
//     from `$CTX_SESSION_ID` when an AI session is active.
//
// First-source-wins ordering means a ref a developer explicitly
// pinned via `ctx trace tag` always shows up before one auto-
// detected from a diff.
//
// # The Trailer
//
// [FormatTrailer] turns a `[]string` of refs into a single git
// trailer line of the form:
//
//	ctx-context: decision:12, task:8, session:abc
//
// Empty input produces an empty string (no trailer is written).
// The trailer is appended to the commit message by the
// `prepare-commit-msg` hook installed by `ctx trace hook
// enable`.
//
// # Persistence
//
// Two append-only JSONL stores live under `state/`:
//
//   - **history.jsonl** — one [HistoryEntry] per commit:
//     full commit hash, the refs that were attached, the
//     commit message, and a UTC timestamp. Written by the
//     `post-commit` hook so the link survives later message
//     edits or squashes.
//   - **overrides.jsonl** — [OverrideEntry] records that let a
//     human pin a different set of refs to a commit after the
//     fact (`ctx trace tag <commit> --note "..."`). Resolution
//     prefers the most recent override over the original
//     history entry.
//
// Both files are read with [ReadHistory] / [ReadOverrides]; both
// silently skip malformed lines so a corrupt tail does not
// break query commands. [WriteHistory] / [WriteOverride] use
// [appendJSONL] which creates the parent directory on demand and
// stamps a UTC timestamp when the caller leaves it zero.
//
// # Resolution
//
// The CLI side (`ctx trace <commit>`, `ctx trace file <path>`)
// asks the package to **resolve** raw refs back to human
// information:
//
//   - [Resolve](ref, contextDir) → [ResolvedRef] with title and
//     one-line preview (or `Found: false` for stale refs).
//   - [CollectRefsForCommit] picks the ref set for a given
//     commit, preferring override → history.
//   - [ResolveCommitHash] takes a short hash, abbrev, or
//     ref-like string and returns the full SHA via `git
//     rev-parse`.
//   - [CommitMessage] / [CommitDate] are thin `git log` wrappers
//     used to render the trace output.
//
// # Concurrency and Safety
//
// All filesystem operations go through [appendJSONL] /
// [readJSONL]; writes are append-only so concurrent commits in
// quick succession (rare but possible with parallel worktrees)
// produce interleaved-but-valid JSONL. The package holds no
// process-wide state.
//
// # Related Packages
//
//   - [github.com/ActiveMemory/ctx/internal/cli/trace] — CLI
//     subcommands (`ctx trace <commit>`, `ctx trace file`,
//     `ctx trace tag`, `ctx trace hook enable/disable`,
//     `ctx trace collect` plumbing).
//   - [github.com/ActiveMemory/ctx/internal/config/trace] —
//     ref-type keywords, file names, trailer format constants.
//   - [github.com/ActiveMemory/ctx/internal/config/dir] — state
//     directory layout (history.jsonl, overrides.jsonl,
//     trace-pending.jsonl all live under `state/`).
//   - [github.com/ActiveMemory/ctx/internal/config/env] —
//     `CTX_SESSION_ID` env var read by [WorkingRefs].
//
// # Background
//
// See `specs/commit-context-tracing.md` and
// `docs/cli/trace.md` for the design rationale and end-user
// documentation.
package trace
