//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package install

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/skill"
	writeSkill "github.com/ActiveMemory/ctx/internal/write/skill"
)

// Cmd returns the "ctx skill install" subcommand.
//
// Returns:
//   - *cobra.Command: Configured install subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySkillInstall)

	return &cobra.Command{
		Use:   cmd.UseSkillInstall,
		Short: short,
		Long:  long,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0])
		},
	}
}

// Run installs a skill from the given source directory.
//
// Parameters:
//   - c: The cobra command for output
//   - source: Path to the source directory containing SKILL.md
func Run(c *cobra.Command, source string) error {
	skillsDir := filepath.Join(rc.ContextDir(), dir.Skills)

	sk, err := skill.Install(source, skillsDir)
	if err != nil {
		return err
	}

	writeSkill.Installed(c, sk.Name, sk.Dir)
	return nil
}
