//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sort provides the file read-order for context assembly.
//
// [ReadOrder] returns context file names in priority order as
// defined by config.FileReadOrder, filtered to files that exist
// in the loaded context. Constitution rules come first, then
// tasks, conventions, architecture, decisions, learnings,
// glossary, and playbook.
package sort
