//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"path/filepath"
	"strings"

	cfgDir "github.com/ActiveMemory/ctx/internal/config/dir"
)

// TrailerKey is the git trailer key used to embed context refs in commit messages.
const TrailerKey = "ctx-context"

// Collect gathers context refs from all three sources — pending records,
// staged file diffs, and current working state — then deduplicates them.
//
// Parameters:
//   - contextDir: absolute path to the .context/ directory
//
// Returns:
//   - []string: deduplicated refs in source order (pending → staged → working)
func Collect(contextDir string) []string {
	stateDir := filepath.Join(contextDir, cfgDir.State)

	var all []string

	// Source 1: pending records written by ctx trace record.
	if entries, err := ReadPending(stateDir); err == nil {
		for _, e := range entries {
			all = append(all, e.Ref)
		}
	}

	// Source 2: staged context file diffs.
	all = append(all, StagedRefs(contextDir)...)

	// Source 3: in-progress tasks and active session env.
	all = append(all, WorkingRefs(contextDir)...)

	return Deduplicate(all)
}

// FormatTrailer formats a slice of refs as a git trailer line.
// Returns an empty string when refs is empty.
//
// Parameters:
//   - refs: context reference strings (e.g. "decision:12", "task:8")
//
// Returns:
//   - string: git trailer like "ctx-context: decision:12, task:8", or ""
func FormatTrailer(refs []string) string {
	if len(refs) == 0 {
		return ""
	}
	return TrailerKey + ": " + strings.Join(refs, ", ")
}

// Deduplicate returns a new slice with duplicate entries removed.
// The first occurrence of each value is preserved; order is maintained.
//
// Parameters:
//   - refs: input slice (may contain duplicates)
//
// Returns:
//   - []string: deduplicated slice, or nil when input is empty
func Deduplicate(refs []string) []string {
	if len(refs) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(refs))
	out := make([]string, 0, len(refs))

	for _, r := range refs {
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
