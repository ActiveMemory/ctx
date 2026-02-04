//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/tpl"
)

// Skills returns the list of embedded skill directory names.
//
// These are Agent Skills (https://agentskills.io) following the specification
// with SKILL.md files containing frontmatter (name, description) and
// autonomy-focused instructions. They can be installed to .claude/skills/
// via "ctx init".
//
// Returns:
//   - []string: Names of available skill directories (e.g., "ctx-status", "ctx-save")
//   - error: Non-nil if the skills directory cannot be read
func Skills() ([]string, error) {
	names, err := tpl.ListSkills()
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}
	return names, nil
}

// SkillContent returns the content of a skill's SKILL.md file by name.
//
// Parameters:
//   - name: Skill directory name as returned by [Skills] (e.g., "ctx-status")
//
// Returns:
//   - []byte: Raw bytes of the SKILL.md file
//   - error: Non-nil if the skill does not exist or cannot be read
func SkillContent(name string) ([]byte, error) {
	content, err := tpl.SkillContent(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill %s: %w", name, err)
	}
	return content, nil
}
