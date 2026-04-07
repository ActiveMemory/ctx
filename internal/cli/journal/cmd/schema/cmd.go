//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/schema/check"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/schema/dump"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the journal schema parent command.
//
// Returns:
//   - *cobra.Command: The schema command with check and dump subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyJournalSchema, cmd.UseJournalSchema,
		check.Cmd(),
		dump.Cmd(),
	)
}
