//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package apply

import (
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/config/token"
	errTask "github.com/ActiveMemory/ctx/internal/err/task"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/task"
)

// completeTask marks a task as complete without output.
//
// Used by [Update] to silently complete tasks detected in the watch
// input stream. Searches for an unchecked task matching the query
// and marks it as done by changing [ ] to [x].
//
// Parameters:
//   - query: search text to match against task descriptions
//     (case-insensitive substring match)
//
// Returns:
//   - error: Non-nil if query is empty, no matching task is found,
//     or file operations fail
func completeTask(query string) error {
	if query == "" {
		return errTask.NoneSpecified()
	}

	filePath := filepath.Join(rc.ContextDir(), ctx.Task)
	nl := token.NewlineLF

	content, readErr := io.SafeReadUserFile(filepath.Clean(filePath))
	if readErr != nil {
		return readErr
	}

	lines := strings.Split(string(content), nl)

	matchedLine := -1
	for i, line := range lines {
		match := regex.Task.FindStringSubmatch(line)
		if match != nil && task.Pending(match) {
			if strings.Contains(
				strings.ToLower(task.Content(match)),
				strings.ToLower(query),
			) {
				matchedLine = i
				break
			}
		}
	}

	if matchedLine == -1 {
		return errTask.NoMatch(query)
	}

	lines[matchedLine] = regex.Task.ReplaceAllString(
		lines[matchedLine], regex.TaskCompleteReplace,
	)
	return io.SafeWriteFile(filePath, []byte(strings.Join(lines, nl)), fs.PermFile)
}
