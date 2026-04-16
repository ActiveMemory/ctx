//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package apply routes **structured context-update
// commands** the AI emitted in its output (XML-tagged
// blocks) to the appropriate per-type writer (task,
// decision, learning, convention, complete).
//
// The package is the dispatch layer behind `ctx watch`,
// which reads stdin (typically piped from `tail -f`
// against an AI's transcript) and applies whatever
// `<ctx-update>` blocks the AI wrote.
//
// # Public Surface
//
//   - **[Update](upd)** — accepts a parsed
//     [entity.PendingUpdate], dispatches to the
//     right backend ([internal/entry] for add,
//     [internal/cli/task] for complete, etc.),
//     and reports the result.
//
// # Update Types
//
// Recognized blocks (each enclosed in
// `<ctx-update type="...">...</ctx-update>`):
//
//   - **task**       — add a new task.
//   - **decision**   — add a new decision.
//   - **learning**   — add a new learning.
//   - **convention** — add a new convention.
//   - **complete**   — mark a task as done.
//
// Unknown types are ignored with a warning so a
// future expansion does not break older installs.
//
// # Idempotency
//
// Add operations are deduplicated by content hash:
// the same block applied twice produces one entry.
// Complete operations match by ID/text and are
// no-ops on already-completed tasks.
//
// # Concurrency
//
// Single-process, sequential within a single
// `ctx watch` invocation.
package apply
