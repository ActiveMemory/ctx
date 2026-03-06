//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"github.com/spf13/cobra"
)

// unpublishCmd returns the memory unpublish subcommand.
//
// Returns:
//   - *cobra.Command: command for removing published context from MEMORY.md.
func unpublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unpublish",
		Short: "Remove published context from MEMORY.md",
		Long: `Remove the ctx-managed marker block from MEMORY.md,
preserving all Claude-owned content outside the markers.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runUnpublish(cmd)
		},
	}
}
