//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// FixHeader prints the "Applying fixes..." heading and a trailing
// blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func FixHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyDriftApplying))
	cmd.Println()
}

// FixRecheck prints the "Re-checking..." message with a leading
// blank line. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
func FixRecheck(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println()
	cmd.Println(desc.Text(text.DescKeyDriftRechecking))
}

// FixedCount prints the number of issues fixed. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of fixed issues
func FixedCount(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftFixedCount), count))
}

// SkippedCount prints the number of issues skipped. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of skipped issues
func SkippedCount(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftSkippedCount), count))
}

// FixError prints a fix error message. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - errMsg: error message string
func FixError(cmd *cobra.Command, errMsg string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftFixError), errMsg))
}

// FixStaleness prints a fix result for a successfully fixed issue.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: file that was fixed
func FixStaleness(cmd *cobra.Command, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftFixStaleness), file))
}

// FixMissing prints a fix result for a created missing file.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: file that was created
func FixMissing(cmd *cobra.Command, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftFixMissing), file))
}

// SkipDeadPath prints a skip message for a dead path reference.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: file containing the dead path
//   - line: line number
//   - path: the dead path
func SkipDeadPath(cmd *cobra.Command, file string, line int, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftSkipDeadPath), file, line, path))
}

// SkipStaleAge prints a skip message for a stale file age warning.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: file with stale age
func SkipStaleAge(cmd *cobra.Command, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftSkipStaleAge), file))
}

// SkipSensitiveFile prints a skip message for a sensitive file violation.
// Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - file: file flagged as sensitive
func SkipSensitiveFile(cmd *cobra.Command, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftSkipSensitiveFile), file))
}

// Archived prints the archive confirmation with task count and
// archive path. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: number of archived tasks
//   - archiveFile: path to the archive file
func Archived(cmd *cobra.Command, count int, archiveFile string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyDriftArchived), count, archiveFile))
}
