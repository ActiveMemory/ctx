//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package change provides terminal output for the change
// log command (ctx change).
//
// # Exported Functions
//
// [List] prints a pre-rendered changes string to stdout.
// The content is assembled by the change core package
// and may include commit history, version diffs, or
// release notes depending on the subcommand invoked.
//
// # Nil Safety
//
// A nil *cobra.Command is treated as a no-op, making
// List safe to call from paths where a command may not
// be available.
//
// # Message Categories
//
//   - Info: rendered change content printed verbatim
//     to stdout without additional formatting
//
// # Usage
//
//	rendered := core.RenderChanges(commits)
//	change.List(cmd, rendered)
package change
