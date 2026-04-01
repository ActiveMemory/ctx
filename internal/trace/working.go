//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	"github.com/ActiveMemory/ctx/internal/task"
)

// envSessionID is the environment variable that carries the active AI session ID.
const envSessionID = "CTX_SESSION_ID"

// WorkingRefs detects context refs from the current working state.
//
// It combines in-progress task refs from TASKS.md with an active AI session
// ref (if CTX_SESSION_ID is set).
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - []string: refs like "task:1", "session:<id>"
func WorkingRefs(contextDir string) []string {
	var refs []string

	refs = append(refs, inProgressTaskRefs(contextDir)...)

	if id := os.Getenv(envSessionID); id != "" {
		refs = append(refs, fmt.Sprintf("session:%s", id))
	}

	return refs
}

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

	//nolint:gosec // path built from trusted contextDir + constant filename
	f, err := os.Open(path)
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
		refs = append(refs, fmt.Sprintf("task:%d", count))
	}

	return refs
}
