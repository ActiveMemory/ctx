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

// ExportPlan prints a dry-run export line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: blob label.
//   - outPath: target file path.
func ExportPlan(cmd *cobra.Command, label, outPath string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyWritePadExportPlan),
			label, outPath,
		),
	)
}

// ExportDone prints a successfully exported blob line.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: blob label.
func ExportDone(cmd *cobra.Command, label string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadExportDone), label))
}

// ErrExportWrite prints a blob write failure to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - label: blob label.
//   - cause: the write error.
func ErrExportWrite(cmd *cobra.Command, label string, cause error) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(
		fmt.Sprintf(
			desc.Text(text.DescKeyWritePadExportWriteFailed),
			label, cause,
		),
	)
}

// ExportSummary prints the export summary or "no blobs" message.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - count: number of blobs exported.
//   - dryRun: whether this was a dry run.
func ExportSummary(cmd *cobra.Command, count int, dryRun bool) {
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
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWritePadExportSummary), verb, count),
	)
}
