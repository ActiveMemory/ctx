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
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	mem "github.com/ActiveMemory/ctx/internal/memory"
	"github.com/ActiveMemory/ctx/internal/rc"
)

func importCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import entries from MEMORY.md into .context/ files",
		Long: `Classify and promote entries from Claude Code's MEMORY.md into
structured .context/ files using heuristic keyword matching.

Each entry is classified as a convention, decision, learning, task,
or skipped (session notes, generic text). Deduplication prevents
re-importing the same entry.

Exit codes:
  0  Imported successfully (or nothing new to import)
  1  MEMORY.md not found`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runImport(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show classification plan without writing")

	return cmd
}

// importResult tracks counts per target for reporting.
type importResult struct {
	conventions int
	decisions   int
	learnings   int
	tasks       int
	skipped     int
	dupes       int
}

func (r importResult) total() int {
	return r.conventions + r.decisions + r.learnings + r.tasks
}

func runImport(cmd *cobra.Command, dryRun bool) error {
	contextDir := rc.ContextDir()
	projectRoot := filepath.Dir(contextDir)

	sourcePath, discoverErr := mem.DiscoverMemoryPath(projectRoot)
	if discoverErr != nil {
		cmd.PrintErrln("Auto memory not active:", discoverErr)
		return fmt.Errorf("MEMORY.md not found")
	}

	sourceData, readErr := os.ReadFile(sourcePath) //nolint:gosec // discovered path
	if readErr != nil {
		return fmt.Errorf("reading MEMORY.md: %w", readErr)
	}

	entries := mem.ParseEntries(string(sourceData))
	if len(entries) == 0 {
		cmd.Println("No entries found in MEMORY.md.")
		return nil
	}

	state, loadErr := mem.LoadState(contextDir)
	if loadErr != nil {
		return fmt.Errorf("loading state: %w", loadErr)
	}

	cmd.Printf("Scanning MEMORY.md for new entries...\n")
	cmd.Printf("  Found %d entries\n\n", len(entries))

	var result importResult

	for _, entry := range entries {
		hash := mem.EntryHash(entry.Text)

		if state.Imported(hash) {
			result.dupes++
			continue
		}

		classification := mem.Classify(entry)
		title := truncate(entry.Text, 60)

		if classification.Target == mem.TargetSkip {
			result.skipped++
			if dryRun {
				cmd.Printf("  -> %q\n     Classified: skip\n\n", title)
			}
			continue
		}

		targetFile := config.FileType[classification.Target]

		if dryRun {
			cmd.Printf("  -> %q\n     Classified: %s (keywords: %s)\n\n",
				title, targetFile, strings.Join(classification.Keywords, ", "))
		} else {
			if promoteErr := mem.Promote(entry, classification); promoteErr != nil {
				cmd.PrintErrf("  Error promoting to %s: %v\n", targetFile, promoteErr)
				continue
			}
			state.MarkImported(hash, classification.Target)
			cmd.Printf("  -> %q\n     Added to %s\n\n", title, targetFile)
		}

		switch classification.Target {
		case config.EntryConvention:
			result.conventions++
		case config.EntryDecision:
			result.decisions++
		case config.EntryLearning:
			result.learnings++
		case config.EntryTask:
			result.tasks++
		}
	}

	// Summary
	if dryRun {
		cmd.Printf("Dry run — would import: %d entries", result.total())
	} else {
		cmd.Printf("Imported: %d entries", result.total())
	}

	var parts []string
	if result.conventions > 0 {
		parts = append(parts, fmt.Sprintf("%d convention", result.conventions))
	}
	if result.decisions > 0 {
		parts = append(parts, fmt.Sprintf("%d decision", result.decisions))
	}
	if result.learnings > 0 {
		parts = append(parts, fmt.Sprintf("%d learning", result.learnings))
	}
	if result.tasks > 0 {
		parts = append(parts, fmt.Sprintf("%d task", result.tasks))
	}
	if len(parts) > 0 {
		cmd.Printf(" (%s)", strings.Join(parts, ", "))
	}
	cmd.Println()

	if result.skipped > 0 {
		cmd.Printf("Skipped: %d entries (session notes/unclassified)\n", result.skipped)
	}
	if result.dupes > 0 {
		cmd.Printf("Duplicates: %d entries (already imported)\n", result.dupes)
	}

	if !dryRun && result.total() > 0 {
		state.MarkImportedDone()
		if saveErr := mem.SaveState(contextDir, state); saveErr != nil {
			return fmt.Errorf("saving state: %w", saveErr)
		}
	}

	return nil
}

func truncate(s string, max int) string {
	// Use first line only
	line := strings.SplitN(s, "\n", 2)[0]
	if len(line) <= max {
		return line
	}
	return line[:max-3] + "..."
}
