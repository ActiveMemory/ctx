//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package consolidate merges consecutive tool runs in journal
// markdown for cleaner reading.
//
// [ToolRuns] detects adjacent tool call/result sections and
// consolidates them into grouped blocks.
package consolidate
