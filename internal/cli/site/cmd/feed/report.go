//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package feed

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/site/core"
)

// printReport outputs the generation summary.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - outPath: Path of the generated feed file
//   - report: Feed generation report with counts and messages
func printReport(cmd *cobra.Command, outPath string, report core.FeedReport) {
	cmd.Println(fmt.Sprintf("\nGenerated %s (%d entries)", outPath, report.Included))

	if len(report.Skipped) > 0 {
		cmd.Println("\nSkipped:")
		for _, msg := range report.Skipped {
			cmd.Println("  " + msg)
		}
	}

	if len(report.Warnings) > 0 {
		cmd.Println("\nWarnings:")
		for _, msg := range report.Warnings {
			cmd.Println("  " + msg)
		}
	}
}
