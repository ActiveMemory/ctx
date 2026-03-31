//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/site/cmd/feed"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx site" parent command.
//
// Subcommands:
//   - feed: Generate an Atom 1.0 feed from blog posts
//
// Returns:
//   - *cobra.Command: Parent command with site subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeySite, cmd.UseSite,
		feed.Cmd(),
	)
}
