//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/obsidian"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/site"
)

// Cmd returns the journal command with subcommands.
//
// The journal system provides LLM-powered analysis and synthesis of
// exported session files in .context/journal/.
//
// Returns:
//   - *cobra.Command: The journal command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyJournal)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyJournal,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(site.Cmd())
	cmd.AddCommand(obsidian.Cmd())

	return cmd
}
