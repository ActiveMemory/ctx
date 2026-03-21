//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mark_journal

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// Cmd returns the "ctx system mark-journal" subcommand.
//
// Returns:
//   - *cobra.Command: Configured mark-journal subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemMarkJournal)

	cmd := &cobra.Command{
		Use:    cmd.UseSystemMarkJournal,
		Short:  short,
		Long:   fmt.Sprintf(long, strings.Join(state.ValidStages, ", ")),
		Hidden: true,
		Args:   cobra.ExactArgs(2), //nolint:mnd // 2 positional args: filename, stage
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMarkJournal(cmd, args[0], args[1])
		},
	}

	cmd.Flags().Bool("check", false, desc.Flag(flag.DescKeySystemMarkJournalCheck))

	return cmd
}
