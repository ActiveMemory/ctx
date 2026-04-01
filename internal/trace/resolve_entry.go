//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"fmt"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/index"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/task"
)

// resolveEntry reads the specified context file, parses entry headers,
// and finds the entry at the given 1-based position.
//
// Parameters:
//   - resolved: partially populated ResolvedRef (Raw, Type, Number already set)
//   - contextDir: absolute path to the .context/ directory
//   - fileName: context file name (e.g. "DECISIONS.md")
//   - number: 1-based entry number
//
// Returns:
//   - ResolvedRef: populated with Title and Detail if found
func resolveEntry(resolved ResolvedRef, contextDir, fileName string, number int) ResolvedRef {
	path := filepath.Clean(filepath.Join(contextDir, fileName))

	content, err := io.SafeReadUserFile(path)
	if err != nil {
		return resolved
	}

	entries := index.ParseHeaders(string(content))

	// Entries are 1-based; index into slice using number-1.
	if number < 1 || number > len(entries) {
		return resolved
	}

	entry := entries[number-1]
	resolved.Title = entry.Title
	resolved.Detail = fmt.Sprintf(desc.Text(text.DescKeyWriteTraceDetailDate), entry.Date)
	resolved.Found = true

	return resolved
}

// resolveTask reads TASKS.md and finds the nth top-level task (1-based),
// counting both pending and completed tasks sequentially.
//
// Parameters:
//   - resolved: partially populated ResolvedRef (Raw, Type, Number already set)
//   - contextDir: absolute path to the .context/ directory
//   - number: 1-based task number
//
// Returns:
//   - ResolvedRef: populated with Title and Detail if found
func resolveTask(resolved ResolvedRef, contextDir string, number int) ResolvedRef {
	path := filepath.Clean(filepath.Join(contextDir, cfgCtx.Task))

	f, err := io.SafeOpenUserFile(path)
	if err != nil {
		return resolved
	}
	defer func() { _ = f.Close() }()

	count := 0
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		m := regex.Task.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Skip subtasks (indented).
		if task.Sub(m) {
			continue
		}

		count++
		if count == number {
			status := cfgTrace.StatusPending
			if task.Completed(m) {
				status = cfgTrace.StatusCompleted
			}
			resolved.Title = task.Content(m)
			resolved.Detail = fmt.Sprintf(
				desc.Text(text.DescKeyWriteTraceDetailStatus), status,
			)
			resolved.Found = true
			return resolved
		}
	}

	return resolved
}
