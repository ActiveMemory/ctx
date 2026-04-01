//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	cfgCtx "github.com/ActiveMemory/ctx/internal/config/ctx"
	"github.com/ActiveMemory/ctx/internal/config/regex"
)

// StagedRefs detects context refs from staged .context/ file diffs.
//
// For each of DECISIONS.md, LEARNINGS.md, and CONVENTIONS.md it runs
// git diff --cached on the file and calls ParseAddedEntries. For TASKS.md
// it calls ParseCompletedTasks. All refs from all files are returned as a
// flat list.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - []string: deduplicated refs found across all staged context files
func StagedRefs(contextDir string) []string {
	type fileEntry struct {
		name      string
		parseFunc func(diff string) []string
	}

	files := []fileEntry{
		{
			name:      cfgCtx.Decision,
			parseFunc: func(diff string) []string { return parseAddedEntries(diff, "decision") },
		},
		{
			name:      cfgCtx.Learning,
			parseFunc: func(diff string) []string { return parseAddedEntries(diff, "learning") },
		},
		{
			name:      cfgCtx.Convention,
			parseFunc: func(diff string) []string { return parseAddedEntries(diff, "convention") },
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

// parseAddedEntries extracts entry refs from added lines in a diff.
//
// Only lines starting with "+" (but not "+++") that match the
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
		// Must start with "+" but not "++" (diff header lines like "+++").
		if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "++") {
			continue
		}
		// Strip the leading "+" before matching.
		content := line[1:]
		if regex.EntryHeader.MatchString(content) {
			count++
			refs = append(refs, fmt.Sprintf("%s:%d", entryType, count))
		}
	}

	if refs == nil {
		return []string{}
	}

	return refs
}

// parseCompletedTasks extracts task refs from newly completed tasks in a diff.
//
// Only lines starting with "+" (but not "+++") that match regex.Task with
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
		// Must start with "+" but not "++" (diff header lines like "+++").
		if !strings.HasPrefix(line, "+") || strings.HasPrefix(line, "++") {
			continue
		}
		// Strip the leading "+" before matching.
		content := line[1:]
		m := regex.Task.FindStringSubmatch(content)
		if m == nil {
			continue
		}
		// m[2] is the state group: "x" for completed, " " or "" for pending.
		if m[2] == "x" {
			count++
			refs = append(refs, fmt.Sprintf("task:%d", count))
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
	//nolint:gosec // filePath is built from trusted contextDir + constant filename by callers
	out, err := exec.Command("git", "diff", "--cached", "--", filePath).Output()
	if err != nil {
		return ""
	}
	return string(out)
}
