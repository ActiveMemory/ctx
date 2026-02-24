//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package permissions implements the "ctx permissions" command for managing
// Claude Code permission snapshots.
//
// The permissions package provides subcommands to:
//   - snapshot: Save settings.local.json as a golden image
//   - restore: Reset settings.local.json from the golden image
//
// Golden images allow teams to maintain a curated permission baseline and
// automatically drop session-accumulated permissions at session start.
package permissions
