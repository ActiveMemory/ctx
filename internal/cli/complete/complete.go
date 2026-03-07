//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"github.com/spf13/cobra"

	completeroot "github.com/ActiveMemory/ctx/internal/cli/complete/cmd/root"
)

// Task finds a task by number or text and marks it complete.
//
// Re-exported from cmd/root for the MCP server, which needs programmatic
// task completion without going through cobra. No other consumer should
// use this — CLI callers go through Cmd().
var Task = completeroot.CompleteTask

// Cmd returns the "ctx complete" command for marking tasks as done.
func Cmd() *cobra.Command {
	return completeroot.Cmd()
}
