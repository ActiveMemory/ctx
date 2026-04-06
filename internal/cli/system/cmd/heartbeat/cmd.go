//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heartbeat

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system heartbeat" subcommand.
//
// Returns:
//   - *cobra.Command: Configured heartbeat subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemHeartbeat)

	return &cobra.Command{
		Use:     cmd.UseSystemHeartbeat,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemHeartbeat),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
