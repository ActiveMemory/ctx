//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package export

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Summary prints what an export will (or would) do based on
// aggregate counters.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - newCount: number of new files to export.
//   - regenCount: number of existing files to regenerate.
//   - skipCount: number of existing files to skip.
//   - lockedCount: number of locked files to skip.
//   - dryRun: when true, uses "Would" instead of "Will".
func Summary(
	cmd *cobra.Command,
	newCount, regenCount, skipCount, lockedCount int,
	dryRun bool,
) {
	if cmd == nil {
		return
	}

	verb := desc.Text(text.DescKeyWriteExportVerb)
	if dryRun {
		verb = desc.Text(text.DescKeyWriteExportVerbDryRun)
	}
	var parts []string
	if newCount > 0 {
		parts = append(parts, fmt.Sprintf(desc.Text(text.DescKeyWriteExportPartNew), newCount))
	}
	if regenCount > 0 {
		parts = append(parts, fmt.Sprintf(desc.Text(text.DescKeyWriteExportPartRegen), regenCount))
	}
	if skipCount > 0 {
		parts = append(parts, fmt.Sprintf(desc.Text(text.DescKeyWriteExportPartSkip), skipCount))
	}
	if lockedCount > 0 {
		parts = append(parts, fmt.Sprintf(desc.Text(text.DescKeyWriteExportPartSkipLock), lockedCount))
	}
	if len(parts) == 0 {
		cmd.Println(desc.Text(text.DescKeyWriteExportNothing))
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteExportSummary), verb, strings.Join(parts, ", ")))
}
