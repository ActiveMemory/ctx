//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prune provides terminal output for the state
// file pruning command (ctx prune).
//
// The prune command removes stale session state files
// that exceed a configured age threshold. Output
// functions cover the full lifecycle from preview
// through completion.
//
// # Dry-Run Preview
//
// [DryRunLine] prints each candidate file with its
// human-readable age so the user can review what
// would be removed before committing.
//
// # Error Reporting
//
// [ErrorLine] writes per-file removal failures to
// stderr. Each line includes the file name and the
// underlying error so the user can investigate.
//
// # Summary
//
// [Summary] closes the operation with counts of
// pruned, skipped, and preserved files. It adjusts
// its wording automatically for dry-run mode,
// showing "would prune" instead of "pruned".
package prune
