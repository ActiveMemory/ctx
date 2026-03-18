//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ctximport

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/spf13/cobra"
)

// NoEntries prints that no entries were found in the source file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name (e.g. "MEMORY.md").
func NoEntries(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportNoEntries), filename))
}

// ScanHeader prints the scanning header: source name, entry count,
// and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name being scanned.
//   - count: number of entries discovered.
func ScanHeader(cmd *cobra.Command, filename string, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportScanning), filename))
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportFound), count))
	cmd.Println()
}

// EntrySkipped prints a skipped entry block: title, "skip"
// classification, and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - title: truncated entry title.
func EntrySkipped(cmd *cobra.Command, title string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportEntry), title))
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteImportClassifiedSkip))
	cmd.Println()
}

// EntryClassified prints a classified entry block (dry run):
// title, target file with keywords, and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - title: truncated entry title.
//   - targetFile: destination filename.
//   - keywords: matched classification keywords.
func EntryClassified(cmd *cobra.Command, title, targetFile string, keywords []string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportEntry), title))
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportClassified), targetFile, strings.Join(keywords, ", ")))
	cmd.Println()
}

// EntryAdded prints a promoted entry block: title, target file,
// and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - title: truncated entry title.
//   - targetFile: destination filename.
func EntryAdded(cmd *cobra.Command, title, targetFile string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportEntry), title))
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportAdded), targetFile))
	cmd.Println()
}

// ErrPromote prints a promotion error to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - targetFile: destination filename.
//   - cause: the promotion error.
func ErrPromote(cmd *cobra.Command, targetFile string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf("  Error promoting to %s: %v", targetFile, cause))
}

// Summary prints the full import summary block: total with
// per-type breakdown, skipped count, and duplicate count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - counts: aggregate import counters.
//   - dryRun: whether this was a dry run.
func Summary(cmd *cobra.Command, result entity.ImportResult, dryRun bool) {
	if cmd == nil {
		return
	}

	total := result.Conventions + result.Decisions + result.Learnings + result.Tasks

	var summary string
	if dryRun {
		summary = fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportSummaryDryRun), total)
	} else {
		summary = fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportSummary), total)
	}

	var parts []string
	if result.Conventions > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountConvention), result.Conventions))
	}
	if result.Decisions > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountDecision), result.Decisions))
	}
	if result.Learnings > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountLearning), result.Learnings))
	}
	if result.Tasks > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountTask), result.Tasks))
	}
	if len(parts) > 0 {
		summary += fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
	}
	cmd.Println(summary)

	if result.Skipped > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportSkipped), result.Skipped))
	}
	if result.Dupes > 0 {
		cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteImportDuplicates), result.Dupes))
	}
}
