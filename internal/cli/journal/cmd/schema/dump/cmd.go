//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dump

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the journal schema dump subcommand.
//
// Returns:
//   - *cobra.Command: Command for printing the embedded schema
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyJournalSchemaDump)

	return &cobra.Command{
		Use:     cmd.UseJournalSchemaDump,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyJournalSchemaDump),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}
