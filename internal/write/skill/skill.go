//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Installed prints confirmation that a skill was installed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
//   - dir: Installation directory
func Installed(cmd *cobra.Command, name, dir string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSkillInstalled),
		name, dir))
}

// NoSkillsFound prints a message indicating no skills are installed.
//
// Parameters:
//   - cmd: Cobra command for output
func NoSkillsFound(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteSkillNoSkills))
}

// EntryWithDesc prints a skill entry with name and description.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
//   - description: Skill description
func EntryWithDesc(cmd *cobra.Command, name, description string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSkillEntryDesc),
		name, description))
}

// Entry prints a skill entry with name only.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
func Entry(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSkillEntry), name))
}

// Count prints the total skill count.
//
// Parameters:
//   - cmd: Cobra command for output
//   - count: Number of skills
func Count(cmd *cobra.Command, count int) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSkillCount), count))
}

// Removed prints confirmation that a skill was removed.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
func Removed(cmd *cobra.Command, name string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSkillRemoved), name))
}
