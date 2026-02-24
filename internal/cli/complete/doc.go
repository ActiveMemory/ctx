//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package complete implements the "ctx complete" command for marking
// tasks as done in TASKS.md.
//
// Tasks can be identified by number or partial text match. The command
// updates TASKS.md by changing "- [ ]" to "- [x]" for the matched task.
package complete
