//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"strings"

	cfgClaude "github.com/ActiveMemory/ctx/internal/config/claude"
)

// skillName extracts the skill name from a permission string like "Skill(name)".
//
// Parameters:
//   - perm: Permission string to parse
//
// Returns:
//   - string: The skill name
//   - bool: True if perm matches the Skill(...) format
func skillName(perm string) (string, bool) {
	if !strings.HasPrefix(perm, cfgClaude.PermSkillPrefix) ||
		!strings.HasSuffix(perm, cfgClaude.PermSkillSuffix) {
		return "", false
	}
	return perm[len(cfgClaude.PermSkillPrefix) : len(perm)-len(cfgClaude.PermSkillSuffix)], true
}
