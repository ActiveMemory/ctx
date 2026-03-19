//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"

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

	short, long := desc.CommandDesc(cmd.DescKeyJournalSite)
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
		&output, "output", "o", defaultOutput, desc.FlagDesc(flag.DescKeyJournalSiteOutput),
	)
	cmd.Flags().BoolVar(
		&build, "build", false, desc.FlagDesc(flag.DescKeyJournalSiteBuild),
	)
	cmd.Flags().BoolVar(
		&serve, "serve", false, desc.FlagDesc(flag.DescKeyJournalSiteServe),
	)

	return cmd
}
