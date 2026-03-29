//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	siteCore "github.com/ActiveMemory/ctx/internal/cli/site/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// PrintFeedReport outputs the feed generation summary.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - outPath: Path of the generated feed file
//   - report: Feed generation report with counts and messages
func PrintFeedReport(
	cmd *cobra.Command,
	outPath string,
	report siteCore.FeedReport,
) {
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
