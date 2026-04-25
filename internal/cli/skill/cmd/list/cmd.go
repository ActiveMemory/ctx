//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

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

// Cmd returns the "ctx skill list" subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeySkillList)

	return &cobra.Command{
		Use:     cmd.UseSkillList,
		Short:   short,
		Example: desc.Example(cmd.DescKeySkillList),
		Args:    cobra.NoArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
}

// Run lists all installed skills with name and description.
//
// Parameters:
//   - c: The cobra command for output
//
// Returns:
//   - error: nil on success, or a skill loading error
func Run(c *cobra.Command) error {
	ctxDir, ctxErr := rc.RequireContextDir()
	if ctxErr != nil {
		c.SilenceUsage = true
		return ctxErr
	}
	skillsDir := filepath.Join(ctxDir, dir.Skills)

	skills, err := skill.LoadAll(skillsDir)
	if err != nil {
		return err
	}

	if len(skills) == 0 {
		writeSkill.NoSkillsFound(c)
		return nil
	}

	for _, sk := range skills {
		if sk.Description != "" {
			writeSkill.EntryWithDesc(c, sk.Name, sk.Description)
		} else {
			writeSkill.Entry(c, sk.Name)
		}
	}

	writeSkill.Count(c, len(skills))
	return nil
}
