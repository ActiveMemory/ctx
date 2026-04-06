//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/ActiveMemory/ctx/internal/flagbind"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the journal obsidian subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating an Obsidian vault from journal
//     entries
func Cmd() *cobra.Command {
	var output string

	short, long := desc.Command(cmd.DescKeyJournalObsidian)
	c := &cobra.Command{
		Use:     cmd.UseJournalObsidian,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalObsidian),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, output)
		},
	}

	defaultOutput := filepath.Join(
		rc.ContextDir(), obsidian.DirName,
	)
	flagbind.StringFlagPDefault(
		c, &output,
		cFlag.Output, cFlag.ShortOutput, defaultOutput,
		flag.DescKeyJournalObsidianOutput,
	)

	return c
}
