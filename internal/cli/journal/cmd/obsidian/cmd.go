//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/obsidian"
	"github.com/spf13/cobra"

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
	cmd := &cobra.Command{
		Use:   "obsidian",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, output)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), obsidian.DirName)
	cmd.Flags().StringVarP(
		&output, "output", "o",
		defaultOutput, desc.FlagDesc(flag.DescKeyJournalObsidianOutput),
	)

	return cmd
}
