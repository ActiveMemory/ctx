//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prune

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/health"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/state"
	"github.com/ActiveMemory/ctx/internal/write/prune"
	"github.com/spf13/cobra"
)

// Run executes the prune logic.
//
// Scans the state directory for session-scoped files (identified by UUID
// patterns) older than the given number of days and removes them. Global
// state files (non-UUID) are preserved. Supports dry-run mode.
//
// Parameters:
//   - cmd: Cobra command for output
//   - days: prune files older than this many days
//   - dryRun: if true, report what would be pruned without removing
//
// Returns:
//   - error: Non-nil on state directory read failure
func Run(cmd *cobra.Command, days int, dryRun bool) error {
	dir := state.StateDir()

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return ctxerr.ReadingDir(readErr)
	}

	cutoff := time.Now().Add(-time.Duration(days) * time2.HoursPerDay * time.Hour)
	var pruned, skipped, preserved int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Only prune files with UUID session IDs
		if !health.UUIDPattern.MatchString(name) {
			preserved++
			continue
		}

		info, statErr := entry.Info()
		if statErr != nil {
			continue
		}

		if info.ModTime().After(cutoff) {
			skipped++
			continue
		}

		if dryRun {
			prune.PruneDryRunLine(cmd, name, health.FormatAge(info.ModTime()))
			pruned++
			continue
		}

		path := filepath.Join(dir, name)
		if rmErr := os.Remove(path); rmErr != nil {
			prune.PruneErrorLine(cmd, name, rmErr)
			continue
		}
		pruned++
	}

	prune.PruneSummary(cmd, dryRun, pruned, skipped, preserved)

	return nil
}
