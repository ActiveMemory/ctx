//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// cleanupMaxAge is the maximum age for temp files before cleanup (15 days).
const cleanupMaxAge = 15 * 24 * time.Hour

// cleanupTmpCmd returns the "ctx system cleanup-tmp" command.
//
// Removes stale files (older than 15 days) from the user-specific ctx
// temp directory on session end.
func cleanupTmpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup-tmp",
		Short: "Clean up stale ctx temp files",
		Long: `Removes files older than 15 days from the user-specific ctx temp
directory ($XDG_RUNTIME_DIR/ctx or /tmp/ctx-<uid>). Runs silently
on session end — no output produced.

Hook event: SessionEnd
Output: none (silent side-effect)`,
		Hidden: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runCleanupTmp()
		},
	}
}

func runCleanupTmp() error {
	tmpDir := secureTempDir()

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil // directory may not exist, not an error
	}

	cutoff := time.Now().Add(-cleanupMaxAge)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(tmpDir, entry.Name()))
		}
	}

	return nil
}
