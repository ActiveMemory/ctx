//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package insert handles **section-aware insertion** of
// new entries into context files, picking the right
// location inside the target file (under the matching
// Phase header, after the latest entry of the same type,
// or at the file bottom as a fallback) instead of just
// appending blindly.
//
// The package is the why-`ctx add` knows-where-to-put-
// things engine. Without it every add would dump at the
// bottom and tasks would lose their phase grouping.
//
// # Public Surface
//
//   - **[AppendEntry](file, entry, opts)**: top-level
//     entry point. Reads `file`, decides where to
//     insert based on `opts.Type` and `opts.Phase`,
//     writes the result back.
//   - **[AfterHeader](lines, header, content)**:
//     pure helper: insert `content` immediately after
//     `header` (or at the end of `header`'s
//     section, depending on the rule). Returns the
//     new line slice.
//   - **[Task](lines, entry, phase)**: task-specific
//     placement: finds the right Phase header (per
//     CONSTITUTION, tasks must stay in their Phase
//     forever) and inserts under it.
//   - **[AppendAtEnd](lines, content)**: fallback
//     when no smarter location can be inferred.
//
// # Constitutional Honors
//
// The TASKS.md rule "tasks stay in their Phase
// section permanently" is enforced here by
// [Task]: a new task always gets the explicit Phase
// header it was added under, never floats free.
//
// # Concurrency
//
// Filesystem-bound. Sequential within a single call.
package insert
