//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package scan

import (
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/file"
	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	execGit "github.com/ActiveMemory/ctx/internal/exec/git"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// FindContextChanges returns context files modified after refTime.
//
// Parameters:
//   - refTime: Only include files modified after this time
//
// Returns:
//   - []entity.ContextChange: Modified files sorted by modtime descending
//   - error: Non-nil if the context directory cannot be read
func FindContextChanges(refTime time.Time) ([]entity.ContextChange, error) {
	dir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return nil, ctxErr
	}
	entries, readDirErr := os.ReadDir(dir)
	if readDirErr != nil {
		return nil, readDirErr
	}

	var changes []entity.ContextChange
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), file.ExtMarkdown) {
			continue
		}
		info, infoErr := e.Info()
		if infoErr != nil {
			continue
		}
		if info.ModTime().After(refTime) {
			changes = append(changes, entity.ContextChange{
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
//   - entity.CodeSummary: Commit count, latest message, dirs, authors
//   - error: Always nil (git failures yield empty summary)
func SummarizeCodeChanges(refTime time.Time) (entity.CodeSummary, error) {
	var summary entity.CodeSummary

	// Count commits.
	out, logErr := GitLogSince(refTime, cfgGit.FlagOneline)
	if logErr != nil {
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
		parts := strings.SplitN(commitLines[0], token.Space, 2)
		if len(parts) == 2 {
			summary.LatestMsg = parts[1]
		}
	}

	// Directories touched.
	dirOut, dirErr := GitLogSince(
		refTime, cfgGit.FlagNameOnly, cfgGit.FormatEmpty, cfgGit.FlagNoCommitID,
	)
	if dirErr == nil {
		summary.Dirs = UniqueTopDirs(string(dirOut))
	}

	// Authors.
	authOut, authErr := GitLogSince(refTime, cfgGit.FormatAuthor)
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
	return execGit.LogSince(t, extraArgs...)
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
	for _, line := range strings.Split(
		strings.TrimSpace(output), token.NewlineLF,
	) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		dir := line
		if i := strings.Index(line, cfgGit.PathSeparator); i >= 0 {
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
	for _, line := range strings.Split(
		strings.TrimSpace(output), token.NewlineLF,
	) {
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
