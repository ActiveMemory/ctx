//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sitereview

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx kb site-review` command.
//
// Returns:
//   - *cobra.Command: configured command.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyKBSiteReview)
	return &cobra.Command{
		Use:   cmd.UseKBSiteReview,
		Short: short,
		Long:  long,
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			return Run(cobraCmd)
		},
	}
}
