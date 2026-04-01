//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package skill provides access to embedded skill directories,
// SKILL.md files, and reference documents.
//
// [List] returns the names of all bundled skills deployed by
// ctx init. [Content] reads a specific skill's SKILL.md file
// by name. Skills are the primary agent instruction mechanism
// in the ctx plugin.
package skill
