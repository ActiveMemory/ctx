//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/memory/core"
	"github.com/ActiveMemory/ctx/internal/config"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// runStatus prints memory bridge status including source location,
// last sync time, line counts, drift indicator, and archive count.
//
// Parameters:
//   - cmd: Cobra command for output routing.
//
// Returns:
//   - error: on discovery failure.
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
			state.LastSync.Local().Format("2006-01-02 15:04"), core.FormatDuration(ago)))
	} else {
		cmd.Println("  Last sync:   never")
	}

	cmd.Println()

	// Source line count
	if sourceData, readErr := os.ReadFile(sourcePath); readErr == nil { //nolint:gosec // discovered path
		line := fmt.Sprintf("  MEMORY.md:  %d lines", core.CountFileLines(sourceData))
		if mem.HasDrift(contextDir, sourcePath) {
			line += " (modified since last sync)"
		}
		cmd.Println(line)
	}

	// Mirror line count
	if mirrorData, readErr := os.ReadFile(mirrorPath); readErr == nil { //nolint:gosec // project-local path
		cmd.Println(fmt.Sprintf("  Mirror:     %d lines", core.CountFileLines(mirrorData)))
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
