//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/build"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/entry"
)

// Cmd returns the "ctx task add" subcommand.
//
// Adds a new task entry to TASKS.md with provenance flags.
// Implementation lives in the shared add core; this thin
// adapter binds the noun and description key.
//
// Returns:
//   - *cobra.Command: Configured task add subcommand
func Cmd() *cobra.Command {
	return build.Cmd(entry.Task, cmd.DescKeyTaskAdd, cmd.UseTaskAdd)
}
