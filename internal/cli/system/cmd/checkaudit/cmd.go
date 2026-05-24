//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkaudit

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system check-audit" subcommand.
//
// Hidden by default: the hook is invoked from
// .claude/settings.local.json (or analogous integration
// config), not by humans.
//
// Returns:
//   - *cobra.Command: Configured check-audit subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySystemCheckAudit)
	return &cobra.Command{
		Use:     cmd.UseSystemCheckAudit,
		Short:   short,
		Example: desc.Example(cmd.DescKeySystemCheckAudit),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
