//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package feed provides the "ctx site feed" subcommand.
package feed

import (
	"path/filepath"

	"github.com/spf13/cobra"
)

// Cmd returns the "ctx site feed" subcommand.
//
// Returns:
//   - *cobra.Command: Configured feed generation subcommand
func Cmd() *cobra.Command {
	var (
		out     string
		baseURL string
	)

	cmd := &cobra.Command{
		Use:   "feed",
		Short: "Generate an Atom 1.0 feed from blog posts",
		Long: `Generate an Atom 1.0 feed from finalized blog posts in docs/blog/.

Parses YAML frontmatter for title, date, author, and topics. Extracts
a summary from the first paragraph after the heading. Only posts with
reviewed_and_finalized: true are included.

Examples:
  ctx site feed
  ctx site feed --out /tmp/feed.xml
  ctx site feed --base-url https://example.com`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runFeed(cmd, "docs/blog", out, baseURL)
		},
	}

	cmd.Flags().StringVarP(
		&out, "out", "o", filepath.Join("site", "feed.xml"),
		"Output path for the generated feed",
	)
	cmd.Flags().StringVar(
		&baseURL, "base-url", "https://ctx.ist",
		"Base URL for entry links",
	)

	return cmd
}
