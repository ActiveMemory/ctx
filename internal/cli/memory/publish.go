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

func publishCmd() *cobra.Command {
	var budget int
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Push curated context to MEMORY.md",
		Long: `Push curated .context/ content into Claude Code's MEMORY.md
so the agent sees structured project context on session start.

Content is wrapped in markers (<!-- ctx:published --> / <!-- ctx:end -->).
Claude-owned content outside the markers is preserved.

Exit codes:
  0  Published successfully
  1  MEMORY.md not found`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runPublish(cmd, budget, dryRun)
		},
	}

	cmd.Flags().IntVar(&budget, "budget", mem.DefaultPublishBudget, "Line budget for published content")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be published without writing")

	return cmd
}

func runPublish(cmd *cobra.Command, budget int, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	memoryPath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		cmd.PrintErrln("Auto memory not active:", discoverErr)
		return fmt.Errorf("MEMORY.md not found")
	}

	result, selectErr := mem.SelectContent(contextDir, budget)
	if selectErr != nil {
		return fmt.Errorf("selecting content: %w", selectErr)
	}

	cmd.Println("Publishing .context/ -> MEMORY.md...")
	cmd.Println()
	cmd.Println("  Source files: TASKS.md, DECISIONS.md, CONVENTIONS.md, LEARNINGS.md")
	cmd.Println(fmt.Sprintf("  Budget: %d lines", budget))
	cmd.Println()
	cmd.Println("  Published block:")
	if len(result.Tasks) > 0 {
		cmd.Println(fmt.Sprintf("    %d pending tasks (from TASKS.md)", len(result.Tasks)))
	}
	if len(result.Decisions) > 0 {
		cmd.Println(fmt.Sprintf("    %d recent decisions (from DECISIONS.md)", len(result.Decisions)))
	}
	if len(result.Conventions) > 0 {
		cmd.Println(fmt.Sprintf("    %d key conventions (from CONVENTIONS.md)", len(result.Conventions)))
	}
	if len(result.Learnings) > 0 {
		cmd.Println(fmt.Sprintf("    %d recent learnings (from LEARNINGS.md)", len(result.Learnings)))
	}
	cmd.Println()
	cmd.Println(fmt.Sprintf("  Total: %d lines (within %d-line budget)", result.TotalLines, budget))

	if dryRun {
		cmd.Println()
		cmd.Println("Dry run — no files written.")
		return nil
	}

	if _, publishErr := mem.Publish(contextDir, memoryPath, budget); publishErr != nil {
		return fmt.Errorf("publishing: %w", publishErr)
	}

	cmd.Println()
	cmd.Println("Published to MEMORY.md (markers: <!-- ctx:published --> ... <!-- ctx:end -->)")

	return nil
}
