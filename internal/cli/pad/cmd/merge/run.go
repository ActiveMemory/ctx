//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/merge"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/store"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/pad"
)

// Run reads entries from input files, deduplicates against the current
// pad, and writes the merged result.
//
// Parameters:
//   - cmd: Cobra command for output
//   - files: Input file paths to merge
//   - keyFile: Optional path to key file (empty = use project key)
//   - dryRun: If true, print summary without writing
//
// Returns:
//   - error: Non-nil on read/write failures
func Run(
	cmd *cobra.Command,
	files []string,
	keyFile string,
	dryRun bool,
) error {
	current, readErr := store.ReadEntries()
	if readErr != nil {
		return readErr
	}

	key := merge.LoadKey(keyFile)

	seen := make(map[string]bool, len(current))
	for _, e := range current {
		seen[e] = true
	}

	blobLabels := merge.BuildBlobLabelMap(current)

	var added, dupes int
	var newEntries []string

	for _, file := range files {
		entries, fileErr := merge.ReadFileEntries(file, key)
		if fileErr != nil {
			return errFs.OpenFile(file, fileErr)
		}

		if merge.HasBinaryEntries(entries) {
			pad.MergeBinaryWarning(cmd, file)
		}

		for _, entry := range entries {
			if seen[entry] {
				dupes++
				pad.MergeDupe(cmd, blob.DisplayEntry(entry))
				continue
			}
			seen[entry] = true

			if conflict, label := merge.HasBlobConflict(entry, blobLabels); conflict {
				pad.MergeBlobConflict(cmd, label)
			}

			newEntries = append(newEntries, entry)
			added++
			pad.MergeAdded(cmd, blob.DisplayEntry(entry), file)
		}
	}

	if added == 0 {
		pad.MergeSummary(cmd, added, dupes, dryRun)
		return nil
	}

	if dryRun {
		pad.MergeSummary(cmd, added, dupes, dryRun)
		return nil
	}

	merged := make([]string, 0, len(current)+len(newEntries))
	merged = append(merged, current...)
	merged = append(merged, newEntries...)
	if writeErr := store.WriteEntries(cmd, merged); writeErr != nil {
		return writeErr
	}

	pad.MergeSummary(cmd, added, dupes, false)
	return nil
}
