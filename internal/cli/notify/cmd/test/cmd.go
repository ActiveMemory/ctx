//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx notify test" subcommand.
//
// Returns:
//   - *cobra.Command: Configured test subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyNotifyTest)
	return &cobra.Command{
		Use:   cmd.UseNotifyTest,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}
