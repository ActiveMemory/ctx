//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package turn handles conversation turn parsing and merging.
//
// [Body] extracts the body text of a conversation turn starting
// from a given line index. [MergeConsecutive] combines adjacent
// turns from the same role into a single block.
package turn
