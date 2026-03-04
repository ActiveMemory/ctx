//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package changes

import (
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/rc"
)

// ContextChange represents a modified context file.
type ContextChange struct {
	Name    string
	ModTime time.Time
}

// CodeSummary summarizes code changes since the reference time.
type CodeSummary struct {
	CommitCount int
	LatestMsg   string
	Dirs        []string
	Authors     []string
}

// FindContextChanges returns context files modified after refTime.
func FindContextChanges(refTime time.Time) ([]ContextChange, error) {
	dir := rc.ContextDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var changes []ContextChange
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		if info.ModTime().After(refTime) {
			changes = append(changes, ContextChange{
				Name:    e.Name(),
				ModTime: info.ModTime(),
			})
		}
	}

	// Sort by modtime descending (most recent first).
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ModTime.After(changes[j].ModTime)
	})

	return changes, nil
}

// SummarizeCodeChanges returns a summary of git activity since refTime.
// All git failures return an empty summary (works in non-git dirs).
func SummarizeCodeChanges(refTime time.Time) (CodeSummary, error) {
	sinceStr := refTime.Format(time.RFC3339)
	var summary CodeSummary

	// Count commits.
	out, err := exec.Command("git", "log", "--oneline", "--since="+sinceStr).Output() //nolint:gosec // time string
	if err != nil {
		return summary, nil
	}
	lines := strings.TrimSpace(string(out))
	if lines == "" {
		return summary, nil
	}
	commitLines := strings.Split(lines, "\n")
	summary.CommitCount = len(commitLines)

	// Latest commit message (first line of oneline output).
	if len(commitLines) > 0 {
		parts := strings.SplitN(commitLines[0], " ", 2)
		if len(parts) == 2 {
			summary.LatestMsg = parts[1]
		}
	}

	// Directories touched.
	dirOut, err := exec.Command("git", "log",
		"--name-only", "--since="+sinceStr,
		"--format=", "--no-commit-id").Output() //nolint:gosec // time string
	if err == nil {
		summary.Dirs = uniqueTopDirs(string(dirOut))
	}

	// Authors.
	authOut, err := exec.Command("git", "log",
		"--since="+sinceStr, "--format=%aN").Output() //nolint:gosec // time string
	if err == nil {
		summary.Authors = uniqueLines(string(authOut))
	}

	return summary, nil
}

// uniqueTopDirs extracts unique top-level directories from file paths.
func uniqueTopDirs(output string) []string {
	seen := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dir := line
		if i := strings.Index(line, "/"); i >= 0 {
			dir = line[:i]
		}
		seen[dir] = true
	}

	dirs := make([]string, 0, len(seen))
	for d := range seen {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)
	return dirs
}

// uniqueLines returns unique non-empty lines from output.
func uniqueLines(output string) []string {
	seen := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			seen[line] = true
		}
	}

	result := make([]string, 0, len(seen))
	for v := range seen {
		result = append(result, v)
	}
	sort.Strings(result)
	return result
}
