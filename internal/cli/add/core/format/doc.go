//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package format renders context entries into their
// Markdown representation for the add command.
//
// # Entry Formatters
//
// Each entry type has a dedicated formatter that returns
// a ready-to-insert Markdown string:
//
//   - [Task] produces a checkbox item with provenance
//     tags. The output includes optional priority, the
//     truncated session ID, git branch, commit hash, and
//     a compact timestamp. Tags use the template strings
//     from the tpl package.
//   - [Decision] produces a structured ADR section with
//     a timestamped heading, status marker, and three
//     subsections: Context, Rationale, and Consequence.
//   - [Learning] produces a structured section with a
//     timestamped heading and three subsections: Context,
//     Lesson, and Application.
//   - [Convention] produces a simple Markdown list item
//     prefixed with "- ".
//
// # Provenance Helpers
//
// Two unexported helpers support the Task formatter:
//
//   - truncateSessionID shortens a UUID to ShortIDLen
//     characters, defaulting to "unknown" when empty.
//   - defaultProvenance returns a value or "unknown"
//     when empty, used for branch and commit fields.
//
// # Data Flow
//
// The cmd/ layer selects a formatter based on the entry
// type predicates in the entry subpackage, passes the
// user-supplied content and metadata, and hands the
// resulting string to insert.AppendEntry for placement
// in the target context file.
package format
