//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package restore provides terminal output for the
// permission restore and snapshot commands
// (ctx permissions restore, ctx permissions snapshot).
//
// These commands manage a golden image of Claude Code
// settings.local.json permissions. Output covers two
// distinct workflows.
//
// # Restore Workflow
//
// [Diff] renders the permission difference between
// the golden image and current settings, showing
// dropped and restored entries for both allow and
// deny rules. When all permission lists are empty
// but the files still differ, it prints a note that
// only non-permission settings changed.
//
// [Done] confirms the restore completed. [NoLocal]
// handles the case where no local settings file
// exists. [Match] reports that the current settings
// already match the golden image.
//
// # Snapshot Workflow
//
// [SnapshotDone] confirms the golden image was
// saved or updated, distinguishing between a new
// save and an update of an existing snapshot.
package restore
