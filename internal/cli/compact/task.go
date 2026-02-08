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
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/context"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// compactTasks moves completed tasks to the "Completed" section in TASKS.md.
//
// Scans TASKS.md for checked items ("- [x]") outside the Completed section,
// including their nested content (indented lines below the task).
// This only moves tasks where all nested subtasks are also complete.
// Optionally archives them to .context/archive/.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - ctx: Loaded context containing the files
//   - archive: If true, write completed tasks to a dated archive file
//
// Returns:
//   - int: Number of tasks moved
//   - error: Non-nil if file write fails
func compactTasks(
	cmd *cobra.Command, ctx *context.Context, archive bool,
) (int, error) {
	var tasksFile *context.FileInfo
	for i := range ctx.Files {
		if ctx.Files[i].Name == config.FileTask {
			tasksFile = &ctx.Files[i]
			break
		}
	}

	if tasksFile == nil {
		return 0, nil
	}

	content := string(tasksFile.Content)
	lines := strings.Split(content, config.NewlineLF)

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Parse task blocks
	blocks := ParseTaskBlocks(lines)

	// Filter to only archivable blocks
	var archivableBlocks []TaskBlock
	for _, block := range blocks {
		if block.IsArchivable {
			archivableBlocks = append(archivableBlocks, block)
			cmd.Println(fmt.Sprintf(
				"%s Moving completed task: %s", green("✓"),
				truncateString(block.ParentTaskText(), 50),
			))
		} else {
			cmd.Println(fmt.Sprintf(
				"%s Skipping (has incomplete children): %s", yellow("!"),
				truncateString(block.ParentTaskText(), 50),
			))
		}
	}

	if len(archivableBlocks) == 0 {
		return 0, nil
	}

	// Remove archivable blocks from lines
	newLines := RemoveBlocksFromLines(lines, archivableBlocks)

	// Add blocks to the Completed section
	for i, line := range newLines {
		if strings.HasPrefix(line, config.HeadingCompleted) {
			// Find the next line that's either empty or another section
			insertIdx := i + 1
			for insertIdx < len(newLines) && newLines[insertIdx] != "" &&
				!strings.HasPrefix(newLines[insertIdx], config.HeadingLevelTwoStart) {
				insertIdx++
			}

			// Build content to insert (full blocks, not just task text)
			var blocksToInsert []string
			for _, block := range archivableBlocks {
				blocksToInsert = append(blocksToInsert, block.Lines...)
			}

			// Insert at the right position
			newLines = append(newLines[:insertIdx],
				append(blocksToInsert, newLines[insertIdx:]...)...,
			)
			break
		}
	}

	// Archive if requested
	if archive && len(archivableBlocks) > 0 {
		// Filter to only tasks old enough to archive
		archiveDays := rc.ArchiveAfterDays()
		var blocksToArchive []TaskBlock
		for _, block := range archivableBlocks {
			if block.OlderThan(archiveDays) {
				blocksToArchive = append(blocksToArchive, block)
			}
		}

		if len(blocksToArchive) > 0 {
			archiveDir := filepath.Join(rc.ContextDir(), config.DirArchive)
			if mkErr := os.MkdirAll(archiveDir, config.PermExec); mkErr == nil {
				archiveFile := filepath.Join(
					archiveDir,
					fmt.Sprintf("tasks-%s.md", time.Now().Format("2006-01-02")),
				)
				nl := config.NewlineLF
				archiveContent := fmt.Sprintf(
					"# Archived Tasks - %s"+nl+nl, time.Now().Format("2006-01-02"),
				)
				for _, block := range blocksToArchive {
					archiveContent += block.BlockContent() + nl + nl
				}
				if writeErr := os.WriteFile(
					archiveFile, []byte(archiveContent), config.PermFile,
				); writeErr == nil {
					cmd.Println(fmt.Sprintf(
						"%s Archived %d tasks to %s (older than %d days)", green("✓"),
						len(blocksToArchive), archiveFile, archiveDays,
					))
				}
			}
		}
	}

	// Write back
	newContent := strings.Join(newLines, config.NewlineLF)
	if newContent != content {
		if err := os.WriteFile(
			tasksFile.Path, []byte(newContent), config.PermFile,
		); err != nil {
			return 0, err
		}
	}

	return len(archivableBlocks), nil
}
