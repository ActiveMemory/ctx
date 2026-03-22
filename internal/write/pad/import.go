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

// ImportNone prints the message when no entries were found to import.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func ImportNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePadImportNone))
}

// ImportDone prints the successful line import count.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of entries imported.
func ImportDone(cmd *cobra.Command, count int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportDone), count))
}

// ImportBlobAdded prints a successfully imported blob line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename of the imported blob.
func ImportBlobAdded(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobAdded), name))
}

// ErrImportBlobSkipped prints a skipped blob to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename.
//   - cause: the error reason.
func ErrImportBlobSkipped(cmd *cobra.Command, name string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobSkipped), name, cause),
	)
}

// ErrImportBlobTooLarge prints a too-large blob skip to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename.
//   - max: maximum allowed size in bytes.
func ErrImportBlobTooLarge(cmd *cobra.Command, name string, max int) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobTooLarge), name, max),
	)
}

// ImportBlobSummary prints the blob import summary or "no files" message.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - added: number of blobs imported.
//   - skipped: number of blobs skipped.
func ImportBlobSummary(cmd *cobra.Command, added, skipped int) {
	if cmd == nil {
		return
	}
	if added == 0 && skipped == 0 {
		cmd.Println(desc.Text(text.DescKeyWritePadImportBlobNone))
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadImportBlobSummary),
			added, skipped,
		),
	)
}

// ErrImportCloseWarning prints a file close warning to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: filename.
//   - cause: the close error.
func ErrImportCloseWarning(cmd *cobra.Command, name string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(
		fmt.Sprintf(
			desc.Text(text.DescKeyWritePadImportCloseWarning), name, cause,
		),
	)
}
