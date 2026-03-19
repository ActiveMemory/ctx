//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package permission implements the "ctx permission" command for managing
// Claude Code permission snapshots.
//
// The permission package provides subcommands to:
//   - snapshot: Save settings.local.json as a golden image
//   - restore: Reset settings.local.json from the golden image
//
// Golden images allow teams to maintain a curated permission baseline and
// automatically drop session-accumulated permissions at session start.
package permission
