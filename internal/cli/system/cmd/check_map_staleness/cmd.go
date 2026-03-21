//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_map_staleness

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system check-map-staleness" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-map-staleness subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemCheckMapStaleness)

	return &cobra.Command{
		Use:    cmd.UseSystemCheckMapStaleness,
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
