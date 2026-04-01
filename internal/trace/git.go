//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"fmt"
	"strings"

	cfgGit "github.com/ActiveMemory/ctx/internal/config/git"
	"github.com/ActiveMemory/ctx/internal/config/token"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/exec/git"
)

// ShortHash returns the first ShortHashLen characters of a commit hash.
//
// Parameters:
//   - hash: full or abbreviated commit hash
//
// Returns:
//   - string: abbreviated hash, or the original if already short
func ShortHash(hash string) string {
	if len(hash) <= cfgTrace.ShortHashLen {
		return hash
	}
	return hash[:cfgTrace.ShortHashLen]
}

// ReadTrailerRefs reads ctx-context trailer refs from a commit message.
//
// Parameters:
//   - commitHash: full commit hash to read trailers from
//
// Returns:
//   - []string: parsed context refs, or empty slice on error
func ReadTrailerRefs(commitHash string) []string {
	out, err := git.Run(
		cfgGit.Log, cfgGit.FlagLast,
		fmt.Sprintf(cfgGit.FormatTrailerValue, cfgTrace.TrailerKey),
		commitHash,
	)
	if err != nil {
		return []string{}
	}

	var refs []string
	for _, line := range strings.Split(
		strings.TrimSpace(string(out)), token.NewlineLF,
	) {
		for _, ref := range strings.Split(
			strings.TrimSpace(line), token.CommaSpace,
		) {
			ref = strings.TrimSpace(ref)
			if ref != "" {
				refs = append(refs, ref)
			}
		}
	}

	return refs
}

// ResolveCommitHash resolves a short ref to a full commit hash.
//
// Parameters:
//   - ref: git commit reference (e.g. "HEAD", "abc1234")
//
// Returns:
//   - string: full commit hash
//   - error: non-nil if git rev-parse fails
func ResolveCommitHash(ref string) (string, error) {
	out, err := git.Run(cfgGit.RevParse, ref)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CommitMessage returns the subject line of a commit.
//
// Parameters:
//   - hash: full commit hash
//
// Returns:
//   - string: commit subject line
//   - error: non-nil if git log fails
func CommitMessage(hash string) (string, error) {
	out, err := git.Run(
		cfgGit.Log, cfgGit.FlagLast, cfgGit.FormatSubject, hash,
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CommitDate returns the commit date string in ISO format.
//
// Parameters:
//   - hash: full commit hash
//
// Returns:
//   - string: commit date, or empty string on error
func CommitDate(hash string) string {
	out, err := git.Run(
		cfgGit.Log, cfgGit.FlagLast, cfgGit.FormatDateISO, hash,
	)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// CollectRefsForCommit gathers context refs for a commit from
// history, overrides, and optionally git trailers.
//
// Parameters:
//   - commitHash: full or abbreviated commit hash
//   - traceDir: absolute path to the trace directory
//   - includeTrailers: whether to read git trailers (slow for bulk)
//
// Returns:
//   - []string: deduplicated refs from all sources
func CollectRefsForCommit(
	commitHash, traceDir string, includeTrailers bool,
) []string {
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
