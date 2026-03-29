//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoSkillsHeader prints the skills list heading.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSkillsHeader(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteSkillsHeader))
	cmd.Println()
}

// InfoSkillLine prints a single skill entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
//   - description: Truncated skill description
func InfoSkillLine(cmd *cobra.Command, name, description string) {
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteSkillLine),
		name, description))
}

// CommandsHeader prints the commands list heading.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func CommandsHeader(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyGuideCommandsHead))
	cmd.Println()
}

// CommandLine prints a single command entry.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: Command name.
//   - short: Short description.
func CommandLine(cmd *cobra.Command, name, short string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyGuideCommandLine), name, short))
}

// Default prints the default guide text.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Default(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Print(desc.Text(text.DescKeyGuideDefault))
}
