//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package change implements the ctx change command, which detects
// context and code changes since the last session or a specified
// time.
//
// It registers the root subcommand and delegates to core/ for
// detection, scanning, and rendering.
package change
