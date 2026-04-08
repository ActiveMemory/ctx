//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rm

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the pad rm subcommand.
//
// Accepts multiple IDs (space-separated) and ranges
// (e.g., "3-5" expands to 3, 4, 5).
//
// Returns:
//   - *cobra.Command: Configured rm subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyPadRm)
	return &cobra.Command{
		Use:     cmd.UsePadRm,
		Short:   short,
		Example: desc.Example(cmd.DescKeyPadRm),
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, parseErr := parse.IDs(args)
			if parseErr != nil {
				return parseErr
			}
			return Run(cmd, ids)
		},
	}
}
