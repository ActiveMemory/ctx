//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package site

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the journal site subcommand.
//
// The --output default is resolved inside [Run] against the
// declared context directory. Computing it at construction time
// would require rc.ContextDir() to succeed before cobra has
// parsed the flags, which is too early under the
// explicit-context-dir model. Leaving the default empty and
// resolving lazily keeps the failure path clean: a missing
// context directory surfaces as a single actionable error from
// Run, not a silently-empty flag default.
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
		Use:     cmd.UseJournalSite,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalSite),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, output, build, serve)
		},
	}

	flagbind.StringFlagPDefault(
		c, &output, cFlag.Output, cFlag.ShortOutput,
		"", flag.DescKeyJournalSiteOutput,
	)
	flagbind.BoolFlag(c, &build, cFlag.Build, flag.DescKeyJournalSiteBuild)
	flagbind.BoolFlag(c, &serve, cFlag.Serve, flag.DescKeyJournalSiteServe)

	return c
}
