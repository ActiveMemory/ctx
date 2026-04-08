//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fmt

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Summary prints the formatting result summary.
//
// Parameters:
//   - cmd: Cobra command for output
//   - formatted: Number of files that were reformatted
//   - total: Total number of context files found
func Summary(cmd *cobra.Command, formatted, total int) {
	cmd.Println(
		fmt.Sprintf(
			desc.Text(text.DescKeyWriteFmtSummary),
			formatted, total,
		),
	)
}

// NeedsFormatting prints a per-file message in check mode.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Context filename that needs formatting
func NeedsFormatting(cmd *cobra.Command, name string) {
	cmd.PrintErrln(
		fmt.Sprintf(
			desc.Text(text.DescKeyWriteFmtNeedsFormatting),
			name,
		),
	)
}
