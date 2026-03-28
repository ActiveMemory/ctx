//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx complete" command for marking tasks as done.
//
// Tasks can be specified by number, partial text match, or full text.
// The command updates TASKS.md by changing "- [ ]" to "- [x]".
//
// Returns:
//   - *cobra.Command: Configured complete command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyComplete)

	cmd := &cobra.Command{
		Use:   "complete <task-id-or-text>",
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE:  Run,
	}

	return cmd
}
