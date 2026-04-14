//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tidy provides the **archive and compact** primitives
// that keep `.context/` files lean as a project ages — moving
// completed tasks into dated archive files, sweeping empty
// sections, and reorganizing TASKS.md without losing
// provenance.
//
// The package is the *engine*; the user-facing surface is
// `ctx task archive`, `ctx compact`, and the `_ctx-archive`
// skill. All three call into the helpers here so the rules are
// applied identically regardless of caller.
//
// # The Archive Pipeline
//
// [WriteArchive](contextDir) is the top-level entry point.
// Behavior:
//
//  1. **Parse** TASKS.md into [TaskBlock] records via
//     [ParseTaskBlocks]. Only **top-level** tasks
//     (`indent == 0`) marked `[x]` are candidates;
//     nested subtasks ride along with their parent.
//  2. **Group** archived tasks by Phase header so the
//     archive file preserves the same Phase structure as
//     the source (a constitutional invariant — Phase
//     identity must survive archival).
//  3. **Write** the archive to
//     `.context/archive/tasks-YYYY-MM-DD.md`, creating
//     the directory if needed. If today's archive file
//     already exists, the new content is *appended*, not
//     overwritten.
//  4. **Remove** the archived blocks from TASKS.md via
//     [RemoveBlocksFromLines] — the rewriter operates on
//     the raw line slice so byte offsets stay aligned.
//
// # Compact and Sanitize
//
// [CompactContext] runs the broader cleanup that `ctx
// compact` performs: archives done tasks **and** sweeps
// empty H2/H3 sections via [RemoveEmptySections] so the
// file does not accumulate dangling headers after every
// archival round. [sanitize.go] holds the helpers that
// trim trailing whitespace, normalize blank-line runs to
// at-most-one, and ensure the file ends with a single
// newline.
//
// # Pure-Logic Core
//
// [block.go] and [parse.go] form the pure-logic core: no
// IO, no time, no flags. They take `[]string` and return
// `[]TaskBlock` / new `[]string`. This split makes the
// archival math testable in isolation; the only IO sits
// in [archive.go] and [compact.go] at the boundary.
//
// # Constitutional Invariants Honored
//
// The CONSTITUTION.md rule "Archival is allowed, deletion
// is not" is enforced at this layer: archival never drops
// content; archive files preserve Phase headers; and
// compaction refuses to touch an entry that has not been
// explicitly marked complete.
//
// # Concurrency
//
// All functions are stateless. Callers serialize through
// process-level execution; concurrent invocations against
// the same context dir would race on file writes (no
// locking is implemented).
//
// # Related Packages
//
//   - [internal/cli/task]            — `ctx task archive`,
//     `ctx task complete`, `ctx task snapshot`.
//   - [internal/cli/compact]         — `ctx compact` CLI
//     entry point.
//   - [internal/mcp/handler]         — MCP `ctx_compact`
//     tool calls into this package.
//   - [internal/entity]              — [TaskBlock] and
//     related domain types.
package tidy
