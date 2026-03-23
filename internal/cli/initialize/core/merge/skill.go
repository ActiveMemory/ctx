//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import "strings"

// skillName extracts the skill name from a permission string like "Skill(name)".
//
// Parameters:
//   - perm: Permission string to parse
//
// Returns:
//   - string: The skill name
//   - bool: True if perm matches the Skill(...) format
func skillName(perm string) (string, bool) {
	if !strings.HasPrefix(perm, "Skill(") || !strings.HasSuffix(perm, ")") {
		return "", false
	}
	return perm[len("Skill(") : len(perm)-1], true
}
