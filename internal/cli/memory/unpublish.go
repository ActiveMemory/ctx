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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func unpublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unpublish",
		Short: "Remove published context from MEMORY.md",
		Long: `Remove the ctx-managed marker block from MEMORY.md,
preserving all Claude-owned content outside the markers.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUnpublish(cmd)
		},
	}
}

func runUnpublish(cmd *cobra.Command) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	memoryPath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		cmd.PrintErrln("Auto memory not active:", discoverErr)
		return fmt.Errorf("MEMORY.md not found")
	}

	data, readErr := os.ReadFile(memoryPath) //nolint:gosec // discovered path
	if readErr != nil {
		return fmt.Errorf("reading MEMORY.md: %w", readErr)
	}

	cleaned, found := mem.RemovePublished(string(data))
	if !found {
		cmd.Println("No published block found in MEMORY.md.")
		return nil
	}

	if writeErr := os.WriteFile(memoryPath, []byte(cleaned), config.PermFile); writeErr != nil {
		return fmt.Errorf("writing MEMORY.md: %w", writeErr)
	}

	cmd.Println("Removed published block from MEMORY.md.")
	return nil
}
