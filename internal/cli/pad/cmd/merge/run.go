//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	ctxerr "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/pad"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core"
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
	current, readErr := core.ReadEntries()
	if readErr != nil {
		return readErr
	}

	key := core.LoadMergeKey(keyFile)

	seen := make(map[string]bool, len(current))
	for _, e := range current {
		seen[e] = true
	}

	blobLabels := core.BuildBlobLabelMap(current)

	var added, dupes int
	var newEntries []string

	for _, file := range files {
		entries, fileErr := core.ReadFileEntries(file, key)
		if fileErr != nil {
			return ctxerr.OpenFile(file, fileErr)
		}

		if core.HasBinaryEntries(entries) {
			pad.PadMergeBinaryWarning(cmd, file)
		}

		for _, entry := range entries {
			if seen[entry] {
				dupes++
				pad.PadMergeDupe(cmd, core.DisplayEntry(entry))
				continue
			}
			seen[entry] = true

			if conflict, label := core.HasBlobConflict(entry, blobLabels); conflict {
				pad.PadMergeBlobConflict(cmd, label)
			}

			newEntries = append(newEntries, entry)
			added++
			pad.PadMergeAdded(cmd, core.DisplayEntry(entry), file)
		}
	}

	if added == 0 {
		pad.PadMergeSummary(cmd, added, dupes, dryRun)
		return nil
	}

	if dryRun {
		pad.PadMergeSummary(cmd, added, dupes, dryRun)
		return nil
	}

	merged := make([]string, 0, len(current)+len(newEntries))
	merged = append(merged, current...)
	merged = append(merged, newEntries...)
	if writeErr := core.WriteEntries(merged); writeErr != nil {
		return writeErr
	}

	pad.PadMergeSummary(cmd, added, dupes, false)
	return nil
}
