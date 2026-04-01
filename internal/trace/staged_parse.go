//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"fmt"
	"strings"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/regex"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/task"
)

// parseAddedEntries extracts entry refs from added lines in a diff.
//
// Only lines starting with "+" (but not "++") that match the
// regex.EntryHeader pattern are counted. Each match produces a ref of
// the form "<entryType>:<count>" where count starts at 1.
//
// Parameters:
//   - diff: output of git diff --cached for a single file
//   - entryType: the type label to use in the returned refs (e.g. "decision")
//
// Returns:
//   - []string: refs like "decision:1", "decision:2"
func parseAddedEntries(diff, entryType string) []string {
	var refs []string
	count := 0

	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, cfgTrace.DiffAddedPrefix) ||
			strings.HasPrefix(line, cfgTrace.DiffHeaderPrefix) {
			continue
		}
		// Strip the leading "+" before matching.
		content := line[1:]
		if regex.EntryHeader.MatchString(content) {
			count++
			refs = append(refs, fmt.Sprintf(cfgTrace.RefFormat, entryType, count))
		}
	}

	if refs == nil {
		return []string{}
	}

	return refs
}

// parseCompletedTasks extracts task refs from newly completed tasks in a diff.
//
// Only lines starting with "+" (but not "++") that match regex.Task with
// state "x" are counted. Each match produces a ref of the form "task:<count>"
// where count starts at 1.
//
// Parameters:
//   - diff: output of git diff --cached for TASKS.md
//
// Returns:
//   - []string: refs like "task:1", "task:2"
func parseCompletedTasks(diff string) []string {
	var refs []string
	count := 0

	scanner := bufio.NewScanner(strings.NewReader(diff))
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, cfgTrace.DiffAddedPrefix) ||
			strings.HasPrefix(line, cfgTrace.DiffHeaderPrefix) {
			continue
		}
		// Strip the leading "+" before matching.
		content := line[1:]
		m := regex.Task.FindStringSubmatch(content)
		if m == nil {
			continue
		}
		if task.Completed(m) {
			count++
			refs = append(refs, fmt.Sprintf(
				cfgTrace.RefFormat, cfgTrace.RefTypeTask, count,
			))
		}
	}

	if refs == nil {
		return []string{}
	}

	return refs
}

// stagedDiff runs git diff --cached -- filePath and returns the output.
// Returns an empty string on any error (best-effort).
//
// Parameters:
//   - filePath: absolute path to the file to diff
//
// Returns:
//   - string: the diff output, or "" on error
func stagedDiff(filePath string) string {
	out, err := git.Run(
		cfgGit.Diff, cfgGit.FlagCached, cfgGit.FlagPathSep, filePath,
	)
	if err != nil {
		return ""
	}
	return string(out)
}
