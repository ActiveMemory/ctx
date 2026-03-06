//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func syncCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Copy MEMORY.md to mirror, archive previous version",
		Long: `Copy Claude Code's MEMORY.md to .context/memory/mirror.md.

Archives the previous mirror before overwriting. Reports line counts
and drift since last sync.

Exit codes:
  0  Synced successfully
  1  MEMORY.md not found (auto memory not active)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSync(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would happen without writing")

	return cmd
}

func runSync(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := memory.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		cmd.PrintErrln("Auto memory not active:", discoverErr)
		return fmt.Errorf("MEMORY.md not found")
	}

	if dryRun {
		cmd.Println("Dry run — no files will be written.")
		cmd.Println(fmt.Sprintf("  Source: %s", sourcePath))
		cmd.Println("  Mirror: .context/memory/mirror.md")
		if memory.HasDrift(contextDir, sourcePath) {
			cmd.Println("  Status: drift detected (source is newer)")
		} else {
			cmd.Println("  Status: no drift")
		}
		return nil
	}

	result, syncErr := memory.Sync(contextDir, sourcePath)
	if syncErr != nil {
		return fmt.Errorf("sync failed: %w", syncErr)
	}

	if result.ArchivedTo != "" {
		cmd.Println(fmt.Sprintf("Archived previous mirror to %s", filepath.Base(result.ArchivedTo)))
	}

	cmd.Println("Synced MEMORY.md -> .context/memory/mirror.md")
	cmd.Println(fmt.Sprintf("  Source: %s", result.SourcePath))
	line := fmt.Sprintf("  Lines: %d", result.SourceLines)
	if result.MirrorLines > 0 {
		line += fmt.Sprintf(" (was %d)", result.MirrorLines)
	}
	cmd.Println(line)

	if result.SourceLines > result.MirrorLines {
		cmd.Println(fmt.Sprintf("  New content: %d lines since last sync",
			result.SourceLines-result.MirrorLines))
	}

	// Update sync state
	state, loadErr := memory.LoadState(contextDir)
	if loadErr != nil {
		return fmt.Errorf("loading state: %w", loadErr)
	}
	state.MarkSynced()
	if saveErr := memory.SaveState(contextDir, state); saveErr != nil {
		return fmt.Errorf("saving state: %w", saveErr)
	}

	return nil
}
