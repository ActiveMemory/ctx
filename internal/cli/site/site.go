//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx site" parent command.
//
// Subcommands:
//   - feed: Generate an Atom 1.0 feed from blog posts
//
// Returns:
//   - *cobra.Command: Parent command with site management subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Site management commands",
		Long: `Manage the ctx.ist static site.

Subcommands:
  feed    Generate an Atom 1.0 feed from blog posts

Examples:
  ctx site feed                              # Generate site/feed.xml
  ctx site feed --out /tmp/feed.xml          # Custom output path
  ctx site feed --base-url https://example.com`,
	}

	cmd.AddCommand(feedCmd())

	return cmd
}
