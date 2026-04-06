//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package remove

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

// Cmd returns the "ctx skill remove" subcommand.
//
// Returns:
//   - *cobra.Command: Configured remove subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySkillRemove)

	return &cobra.Command{
		Use:     cmd.UseSkillRemove,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySkillRemove),
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return Run(c, args[0])
		},
	}
}

// Run removes an installed skill by name.
//
// Parameters:
//   - c: The cobra command for output
//   - name: The skill name to remove
func Run(c *cobra.Command, name string) error {
	skillsDir := filepath.Join(rc.ContextDir(), dir.Skills)

	if err := skill.Remove(skillsDir, name); err != nil {
		return err
	}

	writeSkill.Removed(c, name)
	return nil
}
