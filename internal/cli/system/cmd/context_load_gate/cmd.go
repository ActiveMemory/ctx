//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context_load_gate

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system context-load-gate" subcommand.
//
// Returns:
//   - *cobra.Command: Configured context-load-gate subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemContextLoadGate)

	return &cobra.Command{
		Use:     cmd.UseSystemContextLoadGate,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemContextLoadGate),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
