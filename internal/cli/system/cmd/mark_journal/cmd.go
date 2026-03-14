//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_journal

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// Cmd returns the "ctx system mark-journal" subcommand.
//
// Returns:
//   - *cobra.Command: Configured mark-journal subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemMarkJournal)

	cmd := &cobra.Command{
		Use:    "mark-journal <filename> <stage>",
		Short:  short,
		Long:   fmt.Sprintf(long, strings.Join(state.ValidStages, ", ")),
		Hidden: true,
		Args:   cobra.ExactArgs(2), //nolint:mnd // 2 positional args: filename, stage
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMarkJournal(cmd, args[0], args[1])
		},
	}

	cmd.Flags().Bool("check", false, assets.FlagDesc(assets.FlagDescKeySystemMarkjournalCheck))

	return cmd
}
