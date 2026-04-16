//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sort orders context files by their configured
// read priority for consistent agent context packets.
//
// [ByReadOrder] sorts a slice of [entity.FileInfo] according
// to the [ctx.ReadOrder] configuration. Each file is
// assigned a priority based on its position in the read
// order list: CONSTITUTION.md first, then TASKS.md, and so
// on. Files not in the list receive a fallback priority
// equal to len(ReadOrder), placing them at the end.
//
// The function returns a new sorted slice without modifying
// the original, so callers can sort a working copy while
// preserving the original file order for other uses.
//
// # Why Read Order Matters
//
// The agent context packet presents files in a specific
// sequence so that the most important context (constitution,
// active tasks) appears first. This is critical when the
// agent's context budget is limited: higher-priority files
// are included before the budget is exhausted.
package sort
