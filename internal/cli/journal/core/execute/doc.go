//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package execute runs the import plan, converting session JSONL
// to journal markdown.
//
// [Import] iterates FileActions from the plan, renders each
// session part to Markdown, preserves existing frontmatter when
// regenerating, and writes the output files.
package execute
