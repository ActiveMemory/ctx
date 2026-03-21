//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resources

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system resources" subcommand.
//
// Returns:
//   - *cobra.Command: Configured resources subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemResources)

	cmd := &cobra.Command{
		Use:   "resources",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResources(cmd)
		},
	}
	cmd.Flags().Bool("json", false,
		desc.Flag(flag.DescKeySystemResourcesJson),
	)
	return cmd
}
