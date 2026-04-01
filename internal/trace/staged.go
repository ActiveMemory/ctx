//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"path/filepath"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
)

// StagedRefs detects context refs from staged .context/ file diffs.
//
// For each of DECISIONS.md, LEARNINGS.md, and CONVENTIONS.md it runs
// git diff --cached on the file and calls parseAddedEntries. For TASKS.md
// it calls parseCompletedTasks. All refs from all files are returned as a
// flat list.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - []string: refs found across all staged context files
func StagedRefs(contextDir string) []string {
	type fileEntry struct {
		name      string
		parseFunc func(diff string) []string
	}

	files := []fileEntry{
		{
			name: cfgCtx.Decision,
			parseFunc: func(diff string) []string {
				return parseAddedEntries(diff, cfgTrace.RefTypeDecision)
			},
		},
		{
			name: cfgCtx.Learning,
			parseFunc: func(diff string) []string {
				return parseAddedEntries(diff, cfgTrace.RefTypeLearning)
			},
		},
		{
			name: cfgCtx.Convention,
			parseFunc: func(diff string) []string {
				return parseAddedEntries(diff, cfgTrace.RefTypeConvention)
			},
		},
		{
			name:      cfgCtx.Task,
			parseFunc: parseCompletedTasks,
		},
	}

	var refs []string
	for _, fe := range files {
		path := filepath.Join(contextDir, fe.name)
		diff := stagedDiff(path)
		if diff == "" {
			continue
		}
		refs = append(refs, fe.parseFunc(diff)...)
	}

	if refs == nil {
		return []string{}
	}

	return refs
}
