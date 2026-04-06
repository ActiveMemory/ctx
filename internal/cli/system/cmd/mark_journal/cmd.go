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

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/journal/state"
)

// Cmd returns the "ctx system mark-journal" subcommand.
//
// Returns:
//   - *cobra.Command: Configured mark-journal subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemMarkJournal)

	c := &cobra.Command{
		Use:     cmd.UseSystemMarkJournal,
		Short:   short,
		Long:    fmt.Sprintf(long, strings.Join(state.ValidStages, token.CommaSpace)),
		Example: desc.Example(cmd.DescKeySystemMarkJournal),
		Hidden:  true,
		//nolint:mnd // 2 positional args: filename, stage
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args[0], args[1])
		},
	}

	c.Flags().Bool(
		cFlag.Check,
		false,
		desc.Flag(flag.DescKeySystemMarkJournalCheck),
	)

	return c
}
