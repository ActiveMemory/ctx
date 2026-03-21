//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// PadEmpty prints the message when the scratchpad has no entries.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func PadEmpty(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePadEmpty))
}

// PadKeyCreated prints a key creation notice to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: key file path.
func PadKeyCreated(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(desc.Text(text.DescKeyWritePadKeyCreated), path))
}

// PadEntryAdded prints confirmation that a pad entry was added.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func PadEntryAdded(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryAdded), n))
}

// PadEntryUpdated prints confirmation that a pad entry was updated.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func PadEntryUpdated(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryUpdated), n))
}

// PadExportPlan prints a dry-run export line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: blob label.
//   - outPath: target file path.
func PadExportPlan(cmd *cobra.Command, label, outPath string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadExportPlan), label, outPath))
}

// PadExportDone prints a successfully exported blob line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: blob label.
func PadExportDone(cmd *cobra.Command, label string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadExportDone), label))
}

// ErrPadExportWrite prints a blob write failure to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: blob label.
//   - cause: the write error.
func ErrPadExportWrite(cmd *cobra.Command, label string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(desc.Text(text.DescKeyWritePadExportWriteFailed), label, cause))
}

// PadBlobWritten prints confirmation that a blob was written to a file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - size: number of bytes written.
//   - path: output file path.
func PadBlobWritten(cmd *cobra.Command, size int, path string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadBlobWritten), size, path))
}

// PadEntryRemoved prints confirmation that a pad entry was removed.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func PadEntryRemoved(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryRemoved), n))
}

// PadResolveSide prints a conflict side block: header and numbered entries.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - side: label ("OURS" or "THEIRS").
//   - entries: display strings for each entry.
func PadResolveSide(cmd *cobra.Command, side string, entries []string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadResolveHeader), side))
	for i, entry := range entries {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadResolveEntry), i+1, entry))
	}
}

// PadEntryMoved prints confirmation that a pad entry was moved.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - from: source position (1-based).
//   - to: destination position (1-based).
func PadEntryMoved(cmd *cobra.Command, from, to int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryMoved), from, to))
}

// PadImportNone prints the message when no entries were found to import.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func PadImportNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePadImportNone))
}

// PadImportDone prints the successful line import count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of entries imported.
func PadImportDone(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportDone), count))
}

// PadImportBlobAdded prints a successfully imported blob line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename of the imported blob.
func PadImportBlobAdded(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobAdded), name))
}

// ErrPadImportBlobSkipped prints a skipped blob to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename.
//   - cause: the error reason.
func ErrPadImportBlobSkipped(cmd *cobra.Command, name string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobSkipped), name, cause))
}

// ErrPadImportBlobTooLarge prints a too-large blob skip to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename.
//   - max: maximum allowed size in bytes.
func ErrPadImportBlobTooLarge(cmd *cobra.Command, name string, max int) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobTooLarge), name, max))
}

// PadImportBlobSummary prints the blob import summary or "no files" message.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - added: number of blobs imported.
//   - skipped: number of blobs skipped.
func PadImportBlobSummary(cmd *cobra.Command, added, skipped int) {
	if cmd == nil {
		return
	}
	if added == 0 && skipped == 0 {
		cmd.Println(desc.Text(text.DescKeyWritePadImportBlobNone))
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobSummary), added, skipped))
}

// ErrPadImportCloseWarning prints a file close warning to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename.
//   - cause: the close error.
func ErrPadImportCloseWarning(cmd *cobra.Command, name string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportCloseWarning), name, cause))
}

// PadMergeDupe prints a duplicate-skipped line during merge.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - display: entry display string.
func PadMergeDupe(cmd *cobra.Command, display string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeDupe), display))
}

// PadMergeAdded prints a newly added entry line during merge.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - display: entry display string.
//   - file: source file path.
func PadMergeAdded(cmd *cobra.Command, display, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeAdded), display, file))
}

// PadMergeBlobConflict prints a blob label conflict warning.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: conflicting blob label.
func PadMergeBlobConflict(cmd *cobra.Command, label string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeBlobConflict), label))
}

// PadMergeBinaryWarning prints a binary data warning for a source file.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - file: source file path.
func PadMergeBinaryWarning(cmd *cobra.Command, file string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeBinaryWarning), file))
}

// PadMergeSummary prints the merge summary based on counts and mode.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - added: number of entries added.
//   - dupes: number of duplicates skipped.
//   - dryRun: whether this was a dry run.
func PadMergeSummary(cmd *cobra.Command, added, dupes int, dryRun bool) {
	if cmd == nil {
		return
	}
	if added == 0 && dupes == 0 {
		cmd.Println(desc.Text(text.DescKeyWritePadMergeNone))
		return
	}
	if added == 0 {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeNoneNew), dupes, padPluralize("duplicate", dupes)))
		return
	}
	if dryRun {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeDryRun),
			added, padPluralize("entry", added),
			dupes, padPluralize("duplicate", dupes)))
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeDone),
		added, padPluralize("entry", added),
		dupes, padPluralize("duplicate", dupes)))
}

// padPluralize is an internal helper matching core.Pluralize for write templates.
func padPluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	if len(word) > 0 && word[len(word)-1] == 'y' {
		return word[:len(word)-1] + "ies"
	}
	return word + "s"
}

// PadExportSummary prints the export summary or "no blobs" message.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of blobs exported.
//   - dryRun: whether this was a dry run.
func PadExportSummary(cmd *cobra.Command, count int, dryRun bool) {
	if cmd == nil {
		return
	}
	if count == 0 {
		cmd.Println(desc.Text(text.DescKeyWritePadExportNone))
		return
	}
	verb := desc.Text(text.DescKeyWritePadExportVerbDone)
	if dryRun {
		verb = desc.Text(text.DescKeyWritePadExportVerbDryRun)
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadExportSummary), verb, count))
}
