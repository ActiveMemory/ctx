//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package specs_nudge

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system specs-nudge" subcommand.
//
// Returns:
//   - *cobra.Command: Configured specs-nudge subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemSpecsNudge)

	return &cobra.Command{
		Use:    "specs-nudge",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
