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
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the journal obsidian subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating an Obsidian vault from journal
//     entries
func Cmd() *cobra.Command {
	var output string

	short, long := desc.CommandDesc(cmd.DescKeyJournalObsidian)
	c := &cobra.Command{
		Use:   cmd.UseJournalObsidian,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, output)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), obsidian.DirName)
	c.Flags().StringVarP(
		&output, cFlag.Output, cFlag.ShortOutput,
		defaultOutput, desc.FlagDesc(flag.DescKeyJournalObsidianOutput),
	)

	return c
}
