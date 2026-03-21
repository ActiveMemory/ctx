//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package recall

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/export"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/list"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/lock"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/show"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/sync"
	"github.com/ActiveMemory/ctx/internal/cli/recall/cmd/unlock"
)

// Cmd returns the recall command with subcommands.
//
// The recall system provides commands for browsing and searching AI session
// history across multiple tools (Claude Code, Aider, etc.).
//
// Returns:
//   - *cobra.Command: The recall command with list, show, and serve subcommands
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyRecall)

	cmd := &cobra.Command{
		Use:   cmd.UseRecall,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(list.Cmd())
	cmd.AddCommand(show.Cmd())
	cmd.AddCommand(export.Cmd())
	cmd.AddCommand(lock.Cmd())
	cmd.AddCommand(unlock.Cmd())
	cmd.AddCommand(sync.Cmd())

	return cmd
}
