//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package learnings provides commands for managing LEARNINGS.md.
package learnings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the learnings command with subcommands.
//
// The learnings command provides utilities for managing the LEARNINGS.md file,
// including regenerating the quick-reference index.
//
// Returns:
//   - *cobra.Command: The learnings command with subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "learnings",
		Short: "Manage LEARNINGS.md file",
		Long: `Manage the LEARNINGS.md file and its quick-reference index.

The learnings file maintains an auto-generated index at the top for quick
scanning. Use the subcommands to manage this index.

Subcommands:
  reindex    Regenerate the quick-reference index

Examples:
  ctx learnings reindex`,
	}

	cmd.AddCommand(reindexCmd())

	return cmd
}

// reindexCmd returns the reindex subcommand.
//
// Returns:
//   - *cobra.Command: Command for regenerating the LEARNINGS.md index
func reindexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reindex",
		Short: "Regenerate the quick-reference index",
		Long: `Regenerate the quick-reference index at the top of LEARNINGS.md.

The index is a compact table showing date and title for each learning,
allowing AI agents to quickly scan entries without reading the full file.

This command is useful after manual edits to LEARNINGS.md or when
migrating existing files to use the index format.

Examples:
  ctx learnings reindex`,
		RunE: runReindex,
	}
}

// runReindex regenerates the LEARNINGS.md index.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - args: Command arguments (unused)
//
// Returns:
//   - error: Non-nil if file read/write fails
func runReindex(cmd *cobra.Command, _ []string) error {
	filePath := filepath.Join(rc.GetContextDir(), config.FileLearning)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf(
			"%s not found. Run 'ctx init' first", config.FileLearning,
		)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	updated := index.UpdateLearnings(string(content))

	if err := os.WriteFile(filePath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", filePath, err)
	}

	entries := index.ParseHeaders(string(content))
	green := color.New(color.FgGreen).SprintFunc()
	if len(entries) == 0 {
		cmd.Printf("%s Index cleared (no learnings found)\n", green("✓"))
	} else {
		cmd.Printf(
			"%s Index regenerated with %d entries\n", green("✓"),
			len(entries),
		)
	}

	return nil
}
