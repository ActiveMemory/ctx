//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resources

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system resources" subcommand.
//
// Returns:
//   - *cobra.Command: Configured resources subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemResources)

	cmd := &cobra.Command{
		Use:   "resources",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResources(cmd)
		},
	}
	cmd.Flags().Bool("json", false,
		assets.FlagDesc(assets.FlagDescKeySystemResourcesJson),
	)
	return cmd
}
