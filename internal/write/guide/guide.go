//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"
)

// InfoSkillsHeader prints the skills list heading.
//
// Parameters:
//   - cmd: Cobra command for output
func InfoSkillsHeader(cmd *cobra.Command) {
	cmd.Println(assets.TextDesc(assets.TextDescKeyWriteSkillsHeader))
	cmd.Println()
}

// InfoSkillLine prints a single skill entry.
//
// Parameters:
//   - cmd: Cobra command for output
//   - name: Skill name
//   - description: Truncated skill description
func InfoSkillLine(cmd *cobra.Command, name, description string) {
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWriteSkillLine), name, description))
}
