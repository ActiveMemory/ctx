//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package insert handles section-aware insertion of entries into
// context files.
//
// [AppendEntry] is the main entry point — it reads the target file,
// finds the correct insertion point, and writes the updated content.
// [AfterHeader] inserts below a specific heading, [Task] handles
// task-specific logic (phase sections), and [AppendAtEnd] adds to
// the file bottom as a fallback.
package insert
