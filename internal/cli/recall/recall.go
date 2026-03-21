//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/export"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/lock"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/sync"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/unlock"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the recall command with subcommands.
//
// The recall system provides commands for browsing and searching AI session
// history across multiple tools (Claude Code, Aider, etc.).
//
// Returns:
//   - *cobra.Command: The recall command with list, show, and serve subcommands
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyRecall)

	c := &cobra.Command{
		Use:   cmd.UseRecall,
		Short: short,
		Long:  long,
	}

	c.AddCommand(list.Cmd())
	c.AddCommand(show.Cmd())
	c.AddCommand(export.Cmd())
	c.AddCommand(lock.Cmd())
	c.AddCommand(unlock.Cmd())
	c.AddCommand(sync.Cmd())

	return c
}
