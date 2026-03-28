//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package feed

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/site/core/rss"
	"github.com/ActiveMemory/ctx/internal/cli/site/core/scan"
	writeSite "github.com/ActiveMemory/ctx/internal/write/site"
)

// Run orchestrates scanning and generation.
//
// Parameters:
//   - cmd: Cobra command for output messages
//   - blogDir: Path to the blog posts directory
//   - outPath: Output path for the generated feed
//   - baseURL: Base URL for entry links
//
// Returns:
//   - error: Non-nil if scanning or generation fails
func Run(cmd *cobra.Command, blogDir, outPath, baseURL string) error {
	posts, report, scanErr := scan.BlogPosts(blogDir)
	if scanErr != nil {
		return scanErr
	}

	genErr := rss.Atom(posts, outPath, baseURL)
	if genErr != nil {
		return genErr
	}

	report.Included = len(posts)
	writeSite.PrintFeedReport(cmd, outPath, report)

	return nil
}
