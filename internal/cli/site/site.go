//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/site/cmd/feed"
)

// Cmd returns the "ctx site" parent command.
//
// Subcommands:
//   - feed: Generate an Atom 1.0 feed from blog posts
//
// Returns:
//   - *cobra.Command: Parent command with site management subcommands
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySite)

	cmd := &cobra.Command{
		Use:   "site",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(feed.Cmd())

	return cmd
}
