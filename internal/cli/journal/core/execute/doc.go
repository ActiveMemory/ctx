//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package execute runs the import plan, converting
// session JSONL to journal markdown files.
//
// # Overview
//
// After the user confirms the import plan, this package
// iterates over each file action and writes the
// corresponding journal entry to disk. It handles new
// files, regenerated files (preserving enriched YAML
// frontmatter), and skipped or locked entries.
//
// # Public Surface
//
//   - [Import] -- writes files according to the plan
//     and returns counts of imported, updated, and
//     skipped entries.
//
// # Data Flow
//
// When [Import] is called it processes each FileAction:
//
//  1. Locked actions are skipped with a diagnostic
//     message noting the frontmatter lock.
//  2. Skip actions are skipped with a reason message
//     indicating the file already exists.
//  3. For new and regenerate actions, the session
//     messages are rendered to markdown via the source
//     format sub-package.
//  4. Invalid UTF-8 sequences are replaced with
//     ellipsis characters.
//  5. For regenerated files, the existing YAML
//     frontmatter is preserved unless the discard
//     frontmatter option is set.
//  6. The rendered content is written to disk using
//     safe file I/O.
//  7. The journal state is updated to mark the file
//     as imported.
package execute
