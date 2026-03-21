//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
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

	short, long := desc.Command(cmd.DescKeyJournalSite)
	c := &cobra.Command{
		Use:   cmd.UseJournalSite,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalSite(cmd, output, build, serve)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), dir.JournalSite)
	c.Flags().StringVarP(
		&output, cflag.Output, cflag.ShortOutput, defaultOutput,
		desc.Flag(flag.DescKeyJournalSiteOutput),
	)
	c.Flags().BoolVar(
		&build, cflag.Build, false,
		desc.Flag(flag.DescKeyJournalSiteBuild),
	)
	c.Flags().BoolVar(
		&serve, cflag.Serve, false,
		desc.Flag(flag.DescKeyJournalSiteServe),
	)

	return c
}
