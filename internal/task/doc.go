//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task is the **pure-logic core** behind every
// operation against `TASKS.md` lines: parsing one task
// line into its components, classifying it as completed
// or pending, locating its sub-tasks, and matching it
// against user-supplied selectors.
//
// The package is the foundation of [internal/cli/task]
// (the CLI), [internal/tidy] (the archive engine), and
// [internal/mcp/handler] (the `ctx_complete` MCP tool).
// Everything that touches a TASKS.md line passes through
// the predicates here.
//
// # Public Surface
//
//   - **[Completed](line)** — true when the line is
//     a `- [x] ...` or `- [-] ...` task.
//   - **[Pending](line)** — true when the line is a
//     `- [ ] ...` task.
//   - **[Indent](line)** — returns the leading
//     whitespace count for a task line; used to
//     determine top-level vs nested.
//   - **[Content](line)** — strips the
//     `- [x] `/`- [ ] ` prefix and any trailing
//     inline tags (`#priority:`, `#session:`,
//     `#branch:`, `#commit:`, `#added:`, `#done:`),
//     returning just the human-readable task text.
//   - **[Sub](lines, parentIdx)** — returns the
//     index range of sub-tasks under the task at
//     `parentIdx` (those with strictly greater
//     indent up until the next sibling/parent).
//
// # Why a Separate Package
//
// Five callers need the same predicates and the same
// "what counts as a task line" definition. Hoisting
// them here means the spec lives in one place and the
// audit suite catches duplication.
//
// # Format Reference
//
// Task lines follow the canonical shape established
// by [internal/assets/tpl.Task]:
//
//   - [ ] Implement rate limiting #priority:high
//     #session:abc1 #branch:main #commit:def2
//     #added:2026-04-12-093000
//
// Continuation indents (the wrapped attributes) are
// not separate tasks; [Indent] and the parsers in
// this package treat them as part of the parent
// task's body.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
//
// # Related Packages
//
//   - [internal/cli/task]      — chief consumer for
//     `complete` matching.
//   - [internal/tidy]          — uses the predicates
//     to identify archive candidates.
//   - [internal/mcp/handler]   — uses them in the
//     MCP `ctx_complete` tool.
//   - [internal/assets/tpl]    — defines the task
//     line shape this package parses.
package task
