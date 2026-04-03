//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task moves completed tasks to the Completed section
// and optionally archives them. Tasks with incomplete subtasks
// are skipped. Archive output goes to .context/archive/ when
// enabled via flags or .ctxrc configuration.
package task
