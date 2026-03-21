//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// PrintReport outputs the feed generation summary.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - outPath: Path of the generated feed file
//   - report: Feed generation report with counts and messages
func PrintReport(cmd *cobra.Command, outPath string, report FeedReport) {
	cmd.Println()
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeySiteFeedGenerated), outPath, report.Included))

	if len(report.Skipped) > 0 {
		cmd.Println()
		cmd.Println(desc.Text(text.DescKeySiteFeedSkipped))
		for _, msg := range report.Skipped {
			cmd.Println(fmt.Sprintf(desc.Text(text.DescKeySiteFeedItem), msg))
		}
	}

	if len(report.Warnings) > 0 {
		cmd.Println()
		cmd.Println(desc.Text(text.DescKeySiteFeedWarnings))
		for _, msg := range report.Warnings {
			cmd.Println(fmt.Sprintf(desc.Text(text.DescKeySiteFeedItem), msg))
		}
	}
}
