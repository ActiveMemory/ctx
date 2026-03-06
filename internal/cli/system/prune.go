//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/spf13/cobra"
)

// uuidPattern matches a UUID (v4) anywhere in a filename.
var uuidPattern = regexp.MustCompile(
	`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
)

// pruneCmd returns the "ctx system prune" command.
//
// Session-scoped state files are tombstones and markers that suppress
// repeat hook nudges within a session (e.g. "already checked context
// size", "already sent persistence nudge"). Pruning an active session's
// files causes the corresponding hook to re-fire its nudge — a minor
// UX annoyance, not data loss. No context files, decisions, learnings,
// or code are stored in the state directory.
//
// The one exception is stats-{session}.jsonl, which contains diagnostic
// token usage data. This is informational only and not load-bearing.
func pruneCmd() *cobra.Command {
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

	cmd.Flags().IntVar(&days, "days", 7, "Prune files older than this many days")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be pruned without deleting")

	return cmd
}

func runPrune(cmd *cobra.Command, days int, dryRun bool) error {
	dir := stateDir()

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
		if !uuidPattern.MatchString(name) {
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
			cmd.Println(fmt.Sprintf("  would prune: %s (age: %s)", name, formatAge(info.ModTime())))
			pruned++
			continue
		}

		path := dir + "/" + name
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

func formatAge(t time.Time) string {
	d := time.Since(t)
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
