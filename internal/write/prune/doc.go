//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prune provides terminal output for the state file pruning
// command (ctx system prune).
//
// [DryRunLine] previews each file that would be removed with its
// age. [ErrorLine] reports per-file removal failures. [Summary]
// closes the operation with counts of pruned, skipped, and
// preserved files, adjusting its wording for dry-run mode.
//
// Example:
//
//	for _, f := range stale {
//	    write.DryRunLine(cmd, f.Name, f.Age)
//	}
//	write.Summary(cmd, dryRun, pruned, skipped, preserved)
package prune
