//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package snapshot

import (
	"github.com/spf13/cobra"
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
	cmd := &cobra.Command{
		Use:   "snapshot [name]",
		Short: "Create point-in-time snapshot of TASKS.md",
		Long: `Create a point-in-time snapshot of TASKS.md without modifying the original.

Snapshots are stored in .context/archive/ with timestamped names:
  .context/archive/tasks-snapshot-YYYY-MM-DD-HHMM.md

Unlike archive, snapshot copies the entire file as-is.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runSnapshot,
	}

	return cmd
}
