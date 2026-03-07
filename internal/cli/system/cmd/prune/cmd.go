//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prune

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
)

// Cmd returns the "ctx system prune" subcommand.
//
// Returns:
//   - *cobra.Command: Configured prune subcommand
func Cmd() *cobra.Command {
	var days int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Clean stale per-session state files",
		Long: `Remove per-session state files from .context/state/ that are
older than the specified age. Session state files are identified by
UUID suffixes (e.g. context-check-<session-id>, heartbeat-<session-id>).

Global files without session IDs (events.jsonl, memory-import.json, etc.)
are always preserved.

Examples:
  ctx system prune              # Prune files older than 7 days
  ctx system prune --days 3     # Prune files older than 3 days
  ctx system prune --dry-run    # Show what would be pruned`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPrune(cmd, days, dryRun)
		},
	}

	cmd.Flags().IntVar(&days, "days", 7, assets.FlagDesc(assets.FlagDescKeySystemPruneDays))
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, assets.FlagDesc(assets.FlagDescKeySystemPruneDryRun))

	return cmd
}

func runPrune(cmd *cobra.Command, days int, dryRun bool) error {
	dir := core.StateDir()

	entries, readErr := os.ReadDir(dir)
	if readErr != nil {
		return fmt.Errorf("reading state directory: %w", readErr)
	}

	cutoff := time.Now().Add(-time.Duration(days) * 24 * time.Hour)
	var pruned, skipped, preserved int

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// Only prune files with UUID session IDs
		if !core.UUIDPattern.MatchString(name) {
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
			cmd.Println(fmt.Sprintf("  would prune: %s (age: %s)", name, core.FormatAge(info.ModTime())))
			pruned++
			continue
		}

		path := filepath.Join(dir, name)
		if rmErr := os.Remove(path); rmErr != nil {
			cmd.PrintErrln(fmt.Sprintf("  error removing %s: %v", name, rmErr))
			continue
		}
		pruned++
	}

	if dryRun {
		cmd.Println()
		cmd.Println(fmt.Sprintf("Dry run — would prune %d files (skip %d recent, preserve %d global)",
			pruned, skipped, preserved))
	} else {
		cmd.Println(fmt.Sprintf("Pruned %d files (skipped %d recent, preserved %d global)",
			pruned, skipped, preserved))
	}

	return nil
}
