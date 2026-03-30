//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the tasks snapshot subcommand.
//
// The snapshot command creates a point-in-time copy of TASKS.md without
// modifying the original. Snapshots are stored in .context/archive/ with
// timestamped names.
//
// Arguments:
//   - [name]: Optional name for the snapshot (defaults to "snapshot")
//
// Returns:
//   - *cobra.Command: Configured snapshot subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyTaskSnapshot)

	c := &cobra.Command{
		Use:   cmd.UseTaskSnapshot,
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE:  Run,
	}

	return c
}
