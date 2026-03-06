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

	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func diffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Show what changed since last sync",
		Long: `Show a line-based diff between .context/memory/mirror.md and the
current MEMORY.md. No output when files are identical.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runDiff(cmd)
		},
	}
}

func runDiff(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		return fmt.Errorf("MEMORY.md not found: %w", discoverErr)
	}

	diff, diffErr := mem.Diff(contextDir, sourcePath)
	if diffErr != nil {
		return fmt.Errorf("computing diff: %w", diffErr)
	}

	if diff == "" {
		cmd.Println("No changes since last sync.")
		return nil
	}

	cmd.Print(diff)
	return nil
}
