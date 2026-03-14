//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/token"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
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
//
// Parameters:
//   - refTime: Only include files modified after this time
//
// Returns:
//   - []ContextChange: Modified files sorted by modtime descending
//   - error: Non-nil if the context directory cannot be read
func FindContextChanges(refTime time.Time) ([]ContextChange, error) {
	dir := rc.ContextDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var changes []ContextChange
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), file.ExtMarkdown) {
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
//
// All git failures return an empty summary (works in non-git dirs).
//
// Parameters:
//   - refTime: Only consider commits after this time
//
// Returns:
//   - CodeSummary: Commit count, latest message, dirs, authors
//   - error: Always nil (git failures yield empty summary)
func SummarizeCodeChanges(refTime time.Time) (CodeSummary, error) {
	var summary CodeSummary

	// Count commits.
	out, err := GitLogSince(refTime, "--oneline")
	if err != nil {
		return summary, nil
	}
	lines := strings.TrimSpace(string(out))
	if lines == "" {
		return summary, nil
	}
	commitLines := strings.Split(lines, token.NewlineLF)
	summary.CommitCount = len(commitLines)

	// Latest commit message (first line of oneline output).
	if len(commitLines) > 0 {
		parts := strings.SplitN(commitLines[0], " ", 2)
		if len(parts) == 2 {
			summary.LatestMsg = parts[1]
		}
	}

	// Directories touched.
	dirOut, dirErr := GitLogSince(refTime, "--name-only", "--format=", "--no-commit-id")
	if dirErr == nil {
		summary.Dirs = UniqueTopDirs(string(dirOut))
	}

	// Authors.
	authOut, authErr := GitLogSince(refTime, "--format=%aN")
	if authErr == nil {
		summary.Authors = UniqueLines(string(authOut))
	}

	return summary, nil
}

// GitLogSince runs git log with a --since filter derived from t.
//
// The time value is formatted as RFC 3339 internally so no caller-controlled
// string reaches exec.Command, satisfying gosec G204.
//
// Parameters:
//   - t: Reference time for --since
//   - extraArgs: Additional literal git log flags
//
// Returns:
//   - []byte: Raw git output
//   - error: Non-nil if git fails
func GitLogSince(t time.Time, extraArgs ...string) ([]byte, error) {
	if _, lookErr := exec.LookPath("git"); lookErr != nil {
		return nil, ctxerr.GitNotFound()
	}
	args := []string{"log", "--since", t.Format(time.RFC3339)}
	args = append(args, extraArgs...)
	return exec.Command("git", args...).Output() //nolint:gosec // args are literal flags + time.Format output
}

// UniqueTopDirs extracts unique top-level directories from file paths.
//
// Parameters:
//   - output: Newline-separated file paths
//
// Returns:
//   - []string: Sorted unique top-level directory names
func UniqueTopDirs(output string) []string {
	seen := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(output), token.NewlineLF) {
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

// UniqueLines returns unique non-empty lines from output.
//
// Parameters:
//   - output: Newline-separated text
//
// Returns:
//   - []string: Sorted unique non-empty lines
func UniqueLines(output string) []string {
	seen := make(map[string]bool)
	for _, line := range strings.Split(strings.TrimSpace(output), token.NewlineLF) {
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
