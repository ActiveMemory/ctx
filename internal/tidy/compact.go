//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tidy

import (
	"slices"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
)

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
//   - *compactResult: what changed and what to write
func CompactContext(ctx *entity.Context) *CompactResult {
	result := &CompactResult{}

	// Process TASKS.md.
	tasksFile := ctx.File(cfgCtx.Task)
	if tasksFile != nil {
		content := string(tasksFile.Content)
		lines := strings.Split(content, token.NewlineLF)

		blocks := ParseTaskBlocks(lines)

		var archivableBlocks []entity.TaskBlock
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
				if strings.HasPrefix(line, desc.Text(text.DescKeyHeadingCompleted)) {
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
				result.TasksFileUpdate = &fileUpdate{
					Path:    tasksFile.Path,
					Content: []byte(newContent),
				}
			}
		}
	}

	// Process other files for empty sections.
	for _, f := range ctx.Files {
		if f.Name == cfgCtx.Task {
			continue
		}
		cleaned, count := RemoveEmptySections(string(f.Content))
		if count > 0 {
			result.SectionsCleaned = append(result.SectionsCleaned,
				sectionClean{FileName: f.Name, Removed: count})
			result.SectionFileUpdates = append(result.SectionFileUpdates,
				fileUpdate{Path: f.Path, Content: []byte(cleaned)})
		}
	}

	return result
}
