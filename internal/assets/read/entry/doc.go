//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry provides access to entry template files embedded in
// the assets filesystem.
//
// Entry templates are Markdown scaffolds used when adding new
// decisions, learnings, tasks, or conventions via ctx add.
// [List] returns available template names and [ForName] reads
// a specific template by name.
package entry
