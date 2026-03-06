//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show drift, timestamps, and entry counts",
		Long: `Show memory bridge status: source location, last sync time,
line counts, drift indicator, and archive count.

Exit codes:
  0  No drift
  1  MEMORY.md not found
  2  Drift detected (MEMORY.md changed since last sync)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatus(cmd)
		},
	}
}

func runStatus(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		cmd.Println("Memory Bridge Status")
		cmd.Println("  Source: auto memory not active (MEMORY.md not found)")
		return fmt.Errorf("MEMORY.md not found")
	}

	mirrorPath := filepath.Join(contextDir, config.DirMemory, config.FileMemoryMirror)

	cmd.Println("Memory Bridge Status")
	cmd.Println(fmt.Sprintf("  Source:      %s", sourcePath))
	cmd.Println(fmt.Sprintf("  Mirror:      .context/%s/%s", config.DirMemory, config.FileMemoryMirror))

	// Last sync time
	state, _ := mem.LoadState(contextDir)
	if state.LastSync != nil {
		ago := time.Since(*state.LastSync).Truncate(time.Minute)
		cmd.Println(fmt.Sprintf("  Last sync:   %s (%s ago)",
			state.LastSync.Local().Format("2006-01-02 15:04"), formatDuration(ago)))
	} else {
		cmd.Println("  Last sync:   never")
	}

	cmd.Println()

	// Source line count
	if sourceData, readErr := os.ReadFile(sourcePath); readErr == nil { //nolint:gosec // discovered path
		line := fmt.Sprintf("  MEMORY.md:  %d lines", countFileLines(sourceData))
		if mem.HasDrift(contextDir, sourcePath) {
			line += " (modified since last sync)"
		}
		cmd.Println(line)
	}

	// Mirror line count
	if mirrorData, readErr := os.ReadFile(mirrorPath); readErr == nil { //nolint:gosec // project-local path
		cmd.Println(fmt.Sprintf("  Mirror:     %d lines", countFileLines(mirrorData)))
	} else {
		cmd.Println("  Mirror:     not yet synced")
	}

	// Drift
	hasDrift := mem.HasDrift(contextDir, sourcePath)
	if hasDrift {
		cmd.Println("  Drift:      detected (source is newer)")
	} else {
		cmd.Println("  Drift:      none")
	}

	// Archives
	count := mem.ArchiveCount(contextDir)
	cmd.Println(fmt.Sprintf("  Archives:   %d snapshots in .context/%s/", count, config.DirMemoryArchive))

	if hasDrift {
		// Exit code 2 for drift
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true
		os.Exit(2) //nolint:revive // spec-defined exit code
	}

	return nil
}

func countFileLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	count := 0
	for _, b := range data {
		if b == '\n' {
			count++
		}
	}
	return count
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", m)
	}
	h := int(d.Hours())
	if h == 1 {
		return "1 hour"
	}
	if h < 24 {
		return fmt.Sprintf("%d hours", h)
	}
	days := h / 24
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}
