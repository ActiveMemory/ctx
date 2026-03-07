//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the serve command.
//
// Serves a static site by invoking zensical serve on the specified directory.
//
// Returns:
//   - *cobra.Command: The serve command
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("serve")

	cmd := &cobra.Command{
		Use:   "serve [directory]",
		Short: short,
		Long:  long,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return Run(args)
		},
	}

	return cmd
}
