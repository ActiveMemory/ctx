//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task is the pure-logic core behind every
// operation against TASKS.md lines: parsing one task
// line into its components, classifying it as
// completed or pending, measuring its indent, and
// extracting its human-readable content.
//
// # Public Surface
//
//   - [Completed] -- true when the match represents
//     a checked task (- [x] ...).
//   - [Pending] -- true when the match represents
//     an unchecked task (- [ ] ...).
//   - [Indent] -- returns the leading whitespace
//     from a match, used to determine top-level
//     versus nested tasks.
//   - [Content] -- returns the task text from a
//     match, stripping the checkbox prefix.
//   - [Sub] -- reports whether a match represents
//     a subtask (indented 2+ spaces).
//
// All functions operate on the result of
// ItemPattern.FindStringSubmatch, using the match
// index constants [MatchIndent], [MatchState], and
// [MatchContent].
//
// # Why a Separate Package
//
// Five callers need the same predicates and the same
// definition of what counts as a task line. Hoisting
// them here means the spec lives in one place and
// the audit suite catches duplication.
//
// # Format Reference
//
// Task lines follow the canonical shape established
// by [internal/assets/tpl.Task]:
//
//   - [ ] Implement rate limiting #priority:high
//     #session:abc1 #branch:main #added:2026-04-12
//
// Continuation indents are not separate tasks; the
// parsers treat them as belonging to the parent task body.
//
// # Concurrency
//
// All functions are pure. Concurrent callers never
// race.
package task
