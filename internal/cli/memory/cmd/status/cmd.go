//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx memory status" subcommand.
package status

import (
	"github.com/spf13/cobra"
)

// Cmd returns the memory status subcommand.
//
// Returns:
//   - *cobra.Command: command for showing memory bridge status.
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show drift, timestamps, and entry counts",
		Long: `Show memory bridge status: source location, last sync time,
line counts, drift indicator, and archive count.

Exit codes:
  0  No drift
  1  MEMORY.md not found
  2  Drift detected (MEMORY.md changed since last sync)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runStatus(cmd)
		},
	}
}
