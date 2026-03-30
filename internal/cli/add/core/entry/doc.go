//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry provides type predicates for context file entries.
//
// [FileTypeIsTask], [FileTypeIsDecision], and [FileTypeIsLearning]
// check whether a file type string matches the corresponding entry
// kind. Used by the add command to apply type-specific formatting.
package entry
