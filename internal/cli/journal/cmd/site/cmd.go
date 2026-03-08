//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the journal site subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating a static site from journal entries
func Cmd() *cobra.Command {
	var (
		output string
		serve  bool
		build  bool
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyJournalSite)
	cmd := &cobra.Command{
		Use:   "site",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalSite(cmd, output, build, serve)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), "journal-site")
	cmd.Flags().StringVarP(
		&output, "output", "o", defaultOutput, assets.FlagDesc(assets.FlagDescKeyJournalSiteOutput),
	)
	cmd.Flags().BoolVar(
		&build, "build", false, assets.FlagDesc(assets.FlagDescKeyJournalSiteBuild),
	)
	cmd.Flags().BoolVar(
		&serve, "serve", false, assets.FlagDesc(assets.FlagDescKeyJournalSiteServe),
	)

	return cmd
}
