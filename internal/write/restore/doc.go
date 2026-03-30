//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package restore provides terminal output for the permission
// restore command (ctx permissions restore/snapshot).
//
// Output covers two workflows:
//
//   - Restore: [Diff] shows what would change between the golden
//     image and current settings, [Done] confirms the restore,
//     [NoLocal] handles missing settings, [Match] reports no diff.
//   - Snapshot: [SnapshotDone] confirms the golden image was saved.
package restore
