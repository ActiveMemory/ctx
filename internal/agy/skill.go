//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agy

import "github.com/ActiveMemory/ctx/internal/codex"

// SkillManaged reports whether a SKILL.md file is ctx-managed:
// it opens with a frontmatter block that declares `name: <name>`.
//
// Parameters:
//   - content: existing SKILL.md bytes
//   - name: skill directory name (e.g. "ctx-remember")
//
// Returns:
//   - bool: true when the frontmatter names this skill
func SkillManaged(content []byte, name string) bool {
	return codex.SkillManaged(content, name)
}
