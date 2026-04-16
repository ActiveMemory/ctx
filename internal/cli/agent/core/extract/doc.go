//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract pulls structured items from context
// file content for inclusion in agent context packets.
//
// # Bullet Extraction
//
// [BulletItems] parses Markdown content and returns up
// to a caller-specified limit of bullet list items. It
// strips the "- " prefix, skips empty items and lines
// that start with "#" (headers). The regex pattern comes
// from config/regex.BulletItem.
//
// # Checkbox Extraction
//
// [CheckboxItems] extracts text from both checked and
// unchecked Markdown checkbox items ("- [x]" and
// "- [ ]"). It delegates to config/regex.Task for
// pattern matching and task.Content for field extraction.
//
// # Unchecked Task Filtering
//
// [UncheckedTasks] returns only pending tasks (those
// matching "- [ ]") with the checkbox prefix preserved
// for display. It uses regex.TaskMultiline to handle
// multi-line task bodies and task.Pending to filter.
//
// # Context-Aware Helpers
//
// Two convenience functions operate on a loaded Context:
//
//   - [ActiveTasks] extracts unchecked tasks from the
//     TASKS.md file in the context.
//   - [ConstitutionRules] extracts checkbox items from
//     CONSTITUTION.md for inclusion as inviolable rules.
//
// Both return nil when the target file is absent.
//
// # Data Flow
//
// The budget subpackage calls these functions to populate
// individual sections of the context packet. BulletItems
// feeds into decisions, learnings, and conventions.
// ActiveTasks feeds into the tasks section. Constitution
// rules are emitted first as hard constraints.
package extract
