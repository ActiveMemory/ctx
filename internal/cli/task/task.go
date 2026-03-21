//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package task implements the "ctx tasks" command for managing task archival
// and snapshots.
//
// The task package provides subcommands to:
//   - archive: Move completed tasks to timestamped archive files
//   - snapshot: Create point-in-time copies of TASKS.md
//
// Archive files preserve phase structure for traceability, while snapshots
// copy the entire file as-is without modification.
package task

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/archive"
	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/snapshot"
)

// Cmd returns the task command with subcommands.
//
// The task command provides utilities for managing the task lifecycle:
//   - archive: Move completed tasks out of TASKS.md
//   - snapshot: Create point-in-time backup without modification
//
// Returns:
//   - *cobra.Command: Configured task command with subcommands
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyTask)

	cmd := &cobra.Command{
		Use:   cmd.UseTask,
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(archive.Cmd())
	cmd.AddCommand(complete.Cmd())
	cmd.AddCommand(snapshot.Cmd())

	return cmd
}
