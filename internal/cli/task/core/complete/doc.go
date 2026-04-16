//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package complete finds and marks tasks as done in
// TASKS.md. It supports lookup by task number or by
// case-insensitive text search.
//
// # Task Completion
//
// [Complete] reads TASKS.md, scans for pending tasks
// (lines matching "- [ ]"), and finds the one matching
// the query. The query can be:
//
//   - A number string like "3", which matches the 3rd
//     pending task in document order
//   - A text fragment, which matches any pending task
//     whose text contains the fragment (case-insensitive)
//
// When a match is found, the checkbox is changed from
// "- [ ]" to "- [x]" and the file is written back.
// Returns the matched task text and its 1-based number.
//
// # Error Cases
//
// Complete returns an error when:
//
//   - TASKS.md does not exist
//   - The file cannot be read or written
//   - No pending task matches the query
//   - Multiple pending tasks match a text query
//
// The multiple-match error prevents ambiguous completions
// and asks the user to be more specific.
//
// # Context Directory
//
// The contextDir parameter allows callers to override
// the default context directory from rc.ContextDir.
// Pass an empty string to use the default.
package complete
