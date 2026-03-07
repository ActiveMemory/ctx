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

// CompleteTask finds a task and marks it complete. Re-exported from cmd/root.
var CompleteTask = completeroot.CompleteTask

// Cmd returns the "ctx complete" command for marking tasks as done.
func Cmd() *cobra.Command {
	return completeroot.Cmd()
}
