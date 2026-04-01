//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the ctx change command for detecting
// context and code changes since a reference time.
//
// [Cmd] builds the cobra.Command with --since flag. [Run] resolves
// the reference time (from flag, markers, or event log), scans for
// context file changes and git history, and renders a summary.
package root
