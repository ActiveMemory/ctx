//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// ArchiveKnowledgeFile archives old or superseded entries from a knowledge
// file (DECISIONS.md or LEARNINGS.md).
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - fileName: Context file name (e.g., config.FileDecision)
//   - prefix: Archive file prefix (e.g., "decisions")
//   - heading: Archive heading (e.g., config.HeadingArchivedDecisions)
//   - updateFunc: Reindex function (e.g., index.UpdateDecisions)
//   - days: Entries older than this many days are archived
//   - keepRecent: Number of most recent entries to always keep
//   - archiveAll: If true, archive all entries except keepRecent
//   - dryRun: If true, preview without modifying files
//
// Returns:
//   - int: Number of entries archived
//   - error: Non-nil if file operations fail
func ArchiveKnowledgeFile(
	cmd *cobra.Command,
	fileName, prefix, heading string,
	updateFunc func(string) string,
	days, keepRecent int,
	archiveAll, dryRun bool,
) (int, error) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	filePath := filepath.Join(rc.ContextDir(), fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return 0, fmt.Errorf("no %s found", fileName)
	}

	content, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return 0, fmt.Errorf("failed to read %s: %w", fileName, err)
	}

	blocks := index.ParseEntryBlocks(string(content))
	if len(blocks) == 0 {
		cmd.Println(fmt.Sprintf("No entries found in %s.", fileName))
		return 0, nil
	}

	// Determine which entries to archive
	var toArchive []index.EntryBlock

	// The last keepRecent entries are protected
	protectedStart := len(blocks) - keepRecent
	if protectedStart < 0 {
		protectedStart = 0
	}

	for i, b := range blocks {
		if i >= protectedStart {
			// Protected — skip
			continue
		}
		if archiveAll || b.IsSuperseded() || b.OlderThan(days) {
			toArchive = append(toArchive, b)
		}
	}

	if len(toArchive) == 0 {
		cmd.Println(fmt.Sprintf(
			"No entries to archive in %s (%d entries, all within threshold).",
			fileName, len(blocks),
		))
		return 0, nil
	}

	// Dry-run: preview and return
	if dryRun {
		cmd.Println(yellow("Dry run - no files modified"))
		cmd.Println()
		cmd.Println(fmt.Sprintf(
			"Would archive %d of %d entries from %s (keeping %d recent):",
			len(toArchive), len(blocks), fileName, keepRecent,
		))
		for _, b := range toArchive {
			reason := "old"
			if b.IsSuperseded() {
				reason = "superseded"
			}
			cmd.Println(fmt.Sprintf("  - [%s] %s (%s)", b.Entry.Date, b.Entry.Title, reason))
		}
		return 0, nil
	}

	// Build archive content
	nl := config.NewlineLF
	var archiveContent strings.Builder
	for _, b := range toArchive {
		archiveContent.WriteString(b.BlockContent())
		archiveContent.WriteString(nl + nl)
	}

	// Write to archive file
	archiveFilePath, writeErr := WriteArchive(prefix, heading, archiveContent.String())
	if writeErr != nil {
		return 0, writeErr
	}

	// Remove archived blocks from source
	cleaned := index.RemoveEntryBlocks(string(content), toArchive)

	// Reindex the cleaned content
	cleaned = updateFunc(cleaned)

	// Write back
	if err := os.WriteFile(filePath, []byte(cleaned), config.PermFile); err != nil {
		return 0, fmt.Errorf("failed to write %s: %w", fileName, err)
	}

	cmd.Println(fmt.Sprintf(
		"%s Archived %d entries from %s to %s",
		green("✓"), len(toArchive), fileName, archiveFilePath,
	))

	return len(toArchive), nil
}
