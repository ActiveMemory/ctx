//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_skill_discovery

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system check-skill-discovery" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-skill-discovery subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemCheckSkillDiscovery)

	return &cobra.Command{
		Use:     cmd.UseSystemCheckSkillDiscovery,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemCheckSkillDiscovery),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
