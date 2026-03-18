//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tidy

import (
	"slices"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxCfg "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// FileUpdate holds the new content for a context file that changed
// during compaction.
type FileUpdate struct {
	Path    string
	Content []byte
}

// SectionClean records how many empty sections were removed from a file.
type SectionClean struct {
	FileName string
	Removed  int
}

// CompactResult holds the outcome of a CompactContext call.
//
// Callers decide how to report results (CLI prints, MCP returns
// JSON-RPC responses) and how to write files (os.WriteFile, etc.).
type CompactResult struct {
	// TasksMoved lists the parent text of each task moved to Completed.
	TasksMoved []string
	// TasksSkipped lists parent text of completed tasks with pending children.
	TasksSkipped []string
	// TasksFileUpdate is non-nil when TASKS.md content changed.
	TasksFileUpdate *FileUpdate
	// ArchivableBlocks are blocks eligible for archival.
	ArchivableBlocks []TaskBlock
	// SectionsCleaned lists files where empty sections were removed.
	SectionsCleaned []SectionClean
	// SectionFileUpdates holds new content for files with sections removed.
	SectionFileUpdates []FileUpdate
}

// TotalChanges returns the number of items compacted.
func (r *CompactResult) TotalChanges() int {
	total := len(r.TasksMoved)
	for _, sc := range r.SectionsCleaned {
		total += sc.Removed
	}
	return total
}

// CompactContext analyzes a loaded context for compactable items.
//
// This is pure logic with no I/O. It parses task blocks, moves
// completed tasks to the Completed section, and removes empty
// sections from all context files. Callers are responsible for
// writing the returned file updates to disk.
//
// Parameters:
//   - ctx: loaded context
//
// Returns:
//   - *CompactResult: what changed and what to write
func CompactContext(ctx *entity.Context) *CompactResult {
	result := &CompactResult{}

	// Process TASKS.md.
	tasksFile := ctx.File(ctxCfg.Task)
	if tasksFile != nil {
		content := string(tasksFile.Content)
		lines := strings.Split(content, token.NewlineLF)

		blocks := ParseTaskBlocks(lines)

		var archivableBlocks []TaskBlock
		for _, block := range blocks {
			if block.IsArchivable {
				archivableBlocks = append(archivableBlocks, block)
				result.TasksMoved = append(result.TasksMoved,
					block.ParentTaskText())
			} else {
				result.TasksSkipped = append(result.TasksSkipped,
					block.ParentTaskText())
			}
		}

		if len(archivableBlocks) > 0 {
			result.ArchivableBlocks = archivableBlocks
			newLines := RemoveBlocksFromLines(lines, archivableBlocks)

			// Insert into Completed section.
			for i, line := range newLines {
				if strings.HasPrefix(line, assets.HeadingCompleted) {
					insertIdx := i + 1
					for insertIdx < len(newLines) &&
						newLines[insertIdx] != "" &&
						!strings.HasPrefix(
							newLines[insertIdx],
							token.HeadingLevelTwoStart,
						) {
						insertIdx++
					}

					var blocksToInsert []string
					for _, block := range archivableBlocks {
						blocksToInsert = append(
							blocksToInsert, block.Lines...)
					}

					newLines = slices.Insert(
						newLines, insertIdx, blocksToInsert...)
					break
				}
			}

			newContent := strings.Join(newLines, token.NewlineLF)
			if newContent != content {
				result.TasksFileUpdate = &FileUpdate{
					Path:    tasksFile.Path,
					Content: []byte(newContent),
				}
			}
		}
	}

	// Process other files for empty sections.
	for _, f := range ctx.Files {
		if f.Name == ctxCfg.Task {
			continue
		}
		cleaned, count := RemoveEmptySections(string(f.Content))
		if count > 0 {
			result.SectionsCleaned = append(result.SectionsCleaned,
				SectionClean{FileName: f.Name, Removed: count})
			result.SectionFileUpdates = append(result.SectionFileUpdates,
				FileUpdate{Path: f.Path, Content: []byte(cleaned)})
		}
	}

	return result
}
