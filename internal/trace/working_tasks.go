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

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/task"
)

// inProgressTaskRefs reads TASKS.md and returns a ref for each pending
// top-level task. Subtasks (indent >= 2 spaces) and completed tasks are
// skipped. The first pending top-level task becomes "task:1", the second
// "task:2", and so on.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - []string: refs like "task:1", "task:2" (nil on file read failure)
func inProgressTaskRefs(contextDir string) []string {
	path := filepath.Clean(filepath.Join(contextDir, cfgCtx.Task))

	f, err := io.SafeOpenUserFile(path)
	if err != nil {
		return nil
	}
	defer func() { _ = f.Close() }()

	var refs []string
	count := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := regex.Task.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		if task.Sub(m) {
			continue
		}
		if !task.Pending(m) {
			continue
		}
		count++
		refs = append(refs, fmt.Sprintf(cfgTrace.RefFormat, cfgTrace.RefTypeTask, count))
	}

	return refs
}
