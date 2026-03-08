//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/ActiveMemory/ctx/internal/write/io"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// ImportNoEntries prints that no entries were found in the source file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name (e.g. "MEMORY.md").
func ImportNoEntries(cmd *cobra.Command, filename string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplImportNoEntries, filename)
}

// ImportScanHeader prints the scanning header: source name, entry count,
// and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - filename: source file name being scanned.
//   - count: number of entries discovered.
func ImportScanHeader(cmd *cobra.Command, filename string, count int) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplImportScanning, filename)
	io.sprintf(cmd, config.tplImportFound, count)
	cmd.Println()
}

// ImportEntrySkipped prints a skipped entry block: title, "skip"
// classification, and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - title: truncated entry title.
func ImportEntrySkipped(cmd *cobra.Command, title string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplImportEntry, title)
	cmd.Println(config.tplImportClassifiedSkip)
	cmd.Println()
}

// ImportEntryClassified prints a classified entry block (dry run):
// title, target file with keywords, and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - title: truncated entry title.
//   - targetFile: destination filename.
//   - keywords: matched classification keywords.
func ImportEntryClassified(cmd *cobra.Command, title, targetFile string, keywords []string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplImportEntry, title)
	io.sprintf(cmd, config.tplImportClassified, targetFile, strings.Join(keywords, ", "))
	cmd.Println()
}

// ImportEntryAdded prints a promoted entry block: title, target file,
// and a trailing blank line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - title: truncated entry title.
//   - targetFile: destination filename.
func ImportEntryAdded(cmd *cobra.Command, title, targetFile string) {
	if cmd == nil {
		return
	}
	io.sprintf(cmd, config.tplImportEntry, title)
	io.sprintf(cmd, config.tplImportAdded, targetFile)
	cmd.Println()
}

// ErrImportPromote prints a promotion error to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - targetFile: destination filename.
//   - cause: the promotion error.
func ErrImportPromote(cmd *cobra.Command, targetFile string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf("  Error promoting to %s: %v", targetFile, cause))
}

// ImportCounts holds the per-type tallies for import summary output.
type ImportCounts struct {
	Conventions int
	Decisions   int
	Learnings   int
	Tasks       int
	Skipped     int
	Dupes       int
}

// ImportSummary prints the full import summary block: total with
// per-type breakdown, skipped count, and duplicate count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - counts: aggregate import counters.
//   - dryRun: whether this was a dry run.
func ImportSummary(cmd *cobra.Command, counts ImportCounts, dryRun bool) {
	if cmd == nil {
		return
	}

	total := counts.Conventions + counts.Decisions + counts.Learnings + counts.Tasks

	var summary string
	if dryRun {
		summary = fmt.Sprintf(config.tplImportSummaryDryRun, total)
	} else {
		summary = fmt.Sprintf(config.tplImportSummary, total)
	}

	var parts []string
	if counts.Conventions > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountConvention), counts.Conventions))
	}
	if counts.Decisions > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountDecision), counts.Decisions))
	}
	if counts.Learnings > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountLearning), counts.Learnings))
	}
	if counts.Tasks > 0 {
		parts = append(parts, fmt.Sprintf(
			assets.TextDesc(assets.TextDescKeyImportCountTask), counts.Tasks))
	}
	if len(parts) > 0 {
		summary += fmt.Sprintf(" (%s)", strings.Join(parts, ", "))
	}
	cmd.Println(summary)

	if counts.Skipped > 0 {
		io.sprintf(cmd, config.tplImportSkipped, counts.Skipped)
	}
	if counts.Dupes > 0 {
		io.sprintf(cmd, config.tplImportDuplicates, counts.Dupes)
	}
}
