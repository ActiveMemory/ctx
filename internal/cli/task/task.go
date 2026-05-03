//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package task

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/archive"
	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/complete"
	"github.com/ActiveMemory/ctx/internal/cli/task/cmd/snapshot"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the task command with subcommands.
//
// The task command provides utilities for managing the task
// lifecycle:
//   - add: Add a new task entry to TASKS.md
//   - complete: Mark a task as completed
//   - archive: Move completed tasks out of TASKS.md
//   - snapshot: Create point-in-time backup
//
// Returns:
//   - *cobra.Command: Configured task command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyTask, cmd.UseTask,
		add.Cmd(),
		archive.Cmd(),
		complete.Cmd(),
		snapshot.Cmd(),
	)
}
