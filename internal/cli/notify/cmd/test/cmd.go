//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx notify test" subcommand.
//
// Returns:
//   - *cobra.Command: Configured test subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyNotifyTest)
	return &cobra.Command{
		Use:   "test",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTest(cmd)
		},
	}
}
