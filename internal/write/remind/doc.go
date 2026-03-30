//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package remind provides terminal output for the session reminder
// commands (ctx remind add, list, dismiss).
//
// [Added] confirms a new reminder was created. [Item] renders a
// single reminder in the list with its ID and optional trigger
// condition. [Dismissed] confirms removal. [None] handles the
// empty list case. [DismissedAll] reports bulk dismissal.
package remind
