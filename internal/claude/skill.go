//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"github.com/ActiveMemory/ctx/internal/assets"
)

// Skills returns the list of embedded skill directory names.
//
// These are Agent Skills (https://agentskills.io) following the specification
// with SKILL.md files containing frontmatter (name, description) and
// autonomy-focused instructions. They can be installed to .claude/skills/
// via "ctx init".
//
// Returns:
//   - []string: Names of available skill directories
//     (e.g., "ctx-status", "ctx-reflect")
//   - error: Non-nil if the skills directory cannot be read
func Skills() ([]string, error) {
	names, err := assets.ListSkills()
	if err != nil {
		return nil, errSkillList(err)
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
	content, err := assets.SkillContent(name)
	if err != nil {
		return nil, errSkillRead(name, err)
	}
	return content, nil
}
