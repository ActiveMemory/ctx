//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ShortHash returns the first 7 characters of a commit hash.
func ShortHash(hash string) string {
	if len(hash) <= 7 {
		return hash
	}
	return hash[:7]
}

// ReadTrailerRefs reads ctx-context trailer refs from a commit.
func ReadTrailerRefs(commitHash string) []string {
	//nolint:gosec // TrailerKey is a package constant, commitHash from git rev-parse
	out, err := exec.Command(
		"git", "log", "-1",
		fmt.Sprintf("--format=%%(trailers:key=%s,valueonly)", TrailerKey),
		commitHash,
	).Output()
	if err != nil {
		return []string{}
	}

	var refs []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), token.NewlineLF) {
		for _, ref := range strings.Split(strings.TrimSpace(line), ", ") {
			ref = strings.TrimSpace(ref)
			if ref != "" {
				refs = append(refs, ref)
			}
		}
	}

	return refs
}

// ResolveCommitHash resolves a short ref to a full commit hash.
func ResolveCommitHash(ref string) (string, error) {
	//nolint:gosec // ref is a git commit reference from user input, standard git usage
	out, err := exec.Command("git", "rev-parse", ref).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CommitMessage returns the subject line of a commit.
func CommitMessage(hash string) (string, error) {
	//nolint:gosec // hash is a git commit hash, standard git usage
	out, err := exec.Command("git", "log", "-1", "--format=%s", hash).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CommitDate returns the commit date string.
func CommitDate(hash string) string {
	//nolint:gosec // hash is a git commit hash, standard git usage
	out, err := exec.Command("git", "log", "-1", "--format=%ci", hash).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// CollectRefsForCommit gathers context refs for a commit from
// history, overrides, and optionally git trailers.
func CollectRefsForCommit(commitHash, traceDir string, includeTrailers bool) []string {
	var all []string

	// Source 1: history.jsonl
	if entry, ok := ReadHistoryForCommit(commitHash, traceDir); ok {
		all = append(all, entry.Refs...)
	}

	// Source 2: git trailers (optional — slow for bulk operations)
	if includeTrailers {
		all = append(all, ReadTrailerRefs(commitHash)...)
	}

	// Source 3: overrides.jsonl
	all = append(all, ReadOverridesForCommit(commitHash, traceDir)...)

	result := Deduplicate(all)
	if result == nil {
		return []string{}
	}

	return result
}
