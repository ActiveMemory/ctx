//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"os"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/task/core/count"
	"github.com/ActiveMemory/ctx/internal/cli/task/core/path"
	"github.com/ActiveMemory/ctx/internal/config/archive"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errTask "github.com/ActiveMemory/ctx/internal/err/task"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/tidy"
)

// Plan reads TASKS.md, parses task blocks, and identifies archivable
// blocks. Returns a Result with all data needed for dry-run reporting
// or execution. Does not write any files.
//
// Returns:
//   - Result: Parsed archive plan
//   - error: Non-nil if TASKS.md doesn't exist or can't be read
func Plan() (Result, error) {
	tasksPath := path.FilePath()
	nl := token.NewlineLF

	if _, statErr := os.Stat(tasksPath); os.IsNotExist(statErr) {
		return Result{}, errTask.FileNotFound()
	}

	rawContent, readErr := io.SafeReadUserFile(tasksPath)
	if readErr != nil {
		return Result{}, errTask.FileRead(readErr)
	}

	lines := strings.Split(string(rawContent), nl)
	blocks := tidy.ParseTaskBlocks(lines)

	var r Result
	for _, block := range blocks {
		if block.IsArchivable {
			r.Archivable = append(r.Archivable, block)
		} else {
			r.SkippedNames = append(r.SkippedNames, block.ParentTaskText())
		}
	}

	r.PendingCount = count.Pending(lines)

	var archivedContent strings.Builder
	for _, block := range r.Archivable {
		archivedContent.WriteString(block.BlockContent())
		archivedContent.WriteString(nl)
	}
	r.Content = archivedContent.String()

	newLines := tidy.RemoveBlocksFromLines(lines, r.Archivable)
	r.NewTasksBody = strings.Join(newLines, nl)

	return r, nil
}

// Execute writes the archive file and updates TASKS.md.
//
// Parameters:
//   - r: Result from Plan
//
// Returns:
//   - string: Path to the created archive file
//   - error: Non-nil on write failure
func Execute(r Result) (string, error) {
	archivePath, writeErr := tidy.WriteArchive(
		archive.ScopeTasks,
		desc.Text(text.DescKeyHeadingArchivedTasks),
		r.Content,
	)
	if writeErr != nil {
		return "", writeErr
	}

	tasksPath := path.FilePath()
	if updateErr := os.WriteFile(
		tasksPath, []byte(r.NewTasksBody), fs.PermFile,
	); updateErr != nil {
		return "", errTask.FileWrite(updateErr)
	}

	return archivePath, nil
}
