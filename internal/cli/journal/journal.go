//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/obsidian"
	"github.com/ActiveMemory/ctx/internal/cli/journal/cmd/site"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
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
	c := &cobra.Command{
		Use:   cmd.UseJournal,
		Short: short,
		Long:  long,
	}

	c.AddCommand(site.Cmd())
	c.AddCommand(obsidian.Cmd())

	return c
}
