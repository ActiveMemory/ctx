//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// MergeDupe prints a duplicate-skipped line during merge.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - display: entry display string.
func MergeDupe(cmd *cobra.Command, display string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeDupe), display),
	)
}

// MergeAdded prints a newly added entry line during merge.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - display: entry display string.
//   - file: source file path.
func MergeAdded(cmd *cobra.Command, display, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeAdded), display, file),
	)
}

// MergeBlobConflict prints a blob label conflict warning.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: conflicting blob label.
func MergeBlobConflict(cmd *cobra.Command, label string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeBlobConflict), label),
	)
}

// MergeBinaryWarning prints a binary data warning for a source file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - file: source file path.
func MergeBinaryWarning(cmd *cobra.Command, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeBinaryWarning), file),
	)
}

// MergeSummary prints the merge summary based on counts and mode.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - added: number of entries added.
//   - dupes: number of duplicates skipped.
//   - dryRun: whether this was a dry run.
func MergeSummary(cmd *cobra.Command, added, dupes int, dryRun bool) {
	if cmd == nil {
		return
	}
	if added == 0 && dupes == 0 {
		cmd.Println(desc.Text(text.DescKeyWritePadMergeNone))
		return
	}
	if added == 0 {
		cmd.Println(desc.Text(text.DescKeyWritePadMergeNoneNew))
		mergeSkipped(cmd, dupes)
		return
	}
	if dryRun {
		if added == 1 {
			cmd.Println(desc.Text(text.DescKeyWritePadMergeDryRun1Entry))
		} else {
			cmd.Println(
				fmt.Sprintf(
					desc.Text(text.DescKeyWritePadMergeDryRunNEntries), added),
			)
		}
	} else {
		if added == 1 {
			cmd.Println(desc.Text(text.DescKeyWritePadMergeDone1Entry))
		} else {
			cmd.Println(
				fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeDoneNEntries), added),
			)
		}
	}
	if dupes > 0 {
		mergeSkipped(cmd, dupes)
	}
}
