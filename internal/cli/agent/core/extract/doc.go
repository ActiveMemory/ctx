//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract pulls structured items from context file content.
//
// [BulletItems] extracts markdown list items up to a limit.
// [CheckboxItems] extracts task checkboxes. [UncheckedTasks]
// returns only pending tasks. [ActiveTasks] combines unchecked
// tasks from the loaded context. [ConstitutionRules] extracts
// inviolable rules from CONSTITUTION.md.
package extract
