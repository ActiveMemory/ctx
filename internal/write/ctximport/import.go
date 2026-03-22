//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ctximport

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
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
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWriteImportNoEntries), filename),
	)
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
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteImportScanning), filename))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteImportFound), count))
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
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteImportEntry), title))
	cmd.Println(desc.Text(text.DescKeyWriteImportClassifiedSkip))
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
func EntryClassified(
	cmd *cobra.Command, title, targetFile string, keywords []string,
) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteImportEntry), title))
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyWriteImportClassified),
			targetFile, strings.Join(keywords, token.CommaSpace),
		),
	)
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
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteImportEntry), title))
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteImportAdded), targetFile))
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
	cmd.PrintErrln(
		fmt.Sprintf(
			desc.Text(text.DescKeyWriteImportErrorPromote), targetFile, cause,
		),
	)
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
		summary = fmt.Sprintf(desc.Text(text.DescKeyWriteImportSummaryDryRun), total)
	} else {
		summary = fmt.Sprintf(desc.Text(text.DescKeyWriteImportSummary), total)
	}

	var parts []string
	if result.Conventions > 0 {
		parts = append(parts, fmt.Sprintf(
			desc.Text(text.DescKeyImportCountConvention), result.Conventions))
	}
	if result.Decisions > 0 {
		parts = append(parts, fmt.Sprintf(
			desc.Text(text.DescKeyImportCountDecision), result.Decisions))
	}
	if result.Learnings > 0 {
		parts = append(parts, fmt.Sprintf(
			desc.Text(text.DescKeyImportCountLearning), result.Learnings))
	}
	if result.Tasks > 0 {
		parts = append(parts, fmt.Sprintf(
			desc.Text(text.DescKeyImportCountTask), result.Tasks))
	}
	if len(parts) > 0 {
		summary += fmt.Sprintf(
			desc.Text(text.DescKeyWriteImportBreakdown),
			strings.Join(parts, token.CommaSpace),
		)
	}
	cmd.Println(summary)

	if result.Skipped > 0 {
		cmd.Println(
			fmt.Sprintf(desc.Text(text.DescKeyWriteImportSkipped), result.Skipped),
		)
	}
	if result.Dupes > 0 {
		cmd.Println(
			fmt.Sprintf(desc.Text(text.DescKeyWriteImportDuplicates), result.Dupes),
		)
	}
}
