//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system bootstrap" subcommand.
//
// Returns:
//   - *cobra.Command: Configured bootstrap subcommand
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bootstrap",
		Short: "Print context location for AI agents",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runBootstrap(cmd)
		},
	}
	cmd.Flags().Bool("json", false, "Output in JSON format")
	cmd.Flags().BoolP("quiet", "q", false, "Output only the context directory path")
	return cmd
}
