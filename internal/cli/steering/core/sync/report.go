//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/steering"
	writeSteering "github.com/ActiveMemory/ctx/internal/write/steering"
)

// PrintReport outputs the sync report to the command
// output stream.
func PrintReport(
	c *cobra.Command, report steering.SyncReport,
) {
	for _, name := range report.Written {
		writeSteering.SyncWritten(c, name)
	}
	for _, name := range report.Skipped {
		writeSteering.SyncSkipped(c, name)
	}
	for _, syncErr := range report.Errors {
		writeSteering.SyncError(c, syncErr.Error())
	}

	writeSteering.SyncSummary(c,
		len(report.Written),
		len(report.Skipped),
		len(report.Errors),
	)
}
