//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx notify test" subcommand.
//
// Returns:
//   - *cobra.Command: Configured test subcommand
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Send a test notification",
		Long:  `Sends a test notification to the configured webhook and reports the HTTP status.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runTest(cmd)
		},
	}
}
