//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx compact command for archiving
// completed tasks and cleaning up context files.
//
// [Cmd] builds the cobra.Command with --archive flag. [Run] calls
// core.CompactTasks to move completed tasks to the Completed
// section, removes empty sections, and optionally archives to
// .context/archive/.
package root
