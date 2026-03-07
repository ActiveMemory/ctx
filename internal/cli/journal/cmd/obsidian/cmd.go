//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package obsidian

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Cmd returns the journal obsidian subcommand.
//
// Returns:
//   - *cobra.Command: Command for generating an Obsidian vault from journal
//     entries
func Cmd() *cobra.Command {
	var output string

	short, long := assets.CommandDesc("journal.obsidian")
	cmd := &cobra.Command{
		Use:   "obsidian",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runJournalObsidian(cmd, output)
		},
	}

	defaultOutput := filepath.Join(rc.ContextDir(), config.ObsidianDirName)
	cmd.Flags().StringVarP(
		&output, "output", "o",
		defaultOutput, assets.FlagDesc(assets.FlagDescKeyJournalObsidianOutput),
	)

	return cmd
}
