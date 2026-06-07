//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package review

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the dream review subcommand.
//
// Returns:
//   - *cobra.Command: configured review subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyDreamReview)
	return &cobra.Command{
		Use:     cmd.UseDreamReview,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyDreamReview),
		Args:    cobra.NoArgs,
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			return Run(cobraCmd)
		},
	}
}
