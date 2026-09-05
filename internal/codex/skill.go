//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package codex

import (
	"bufio"
	"bytes"
	"strings"

	cfgCodex "github.com/ActiveMemory/ctx/internal/config/codex"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// SkillManaged reports whether a SKILL.md file is ctx-managed:
// it opens with a frontmatter block that declares `name: <name>`.
// A foreign file at a ctx skill path fails this check and must
// not be overwritten.
//
// Parameters:
//   - content: existing SKILL.md bytes
//   - name: skill directory name (e.g. "ctx-remember")
//
// Returns:
//   - bool: true when the frontmatter names this skill
func SkillManaged(content []byte, name string) bool {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if !scanner.Scan() ||
		strings.TrimSpace(scanner.Text()) != token.FrontmatterDelimiter {
		return false
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == token.FrontmatterDelimiter {
			return false
		}
		key, value, hasSep := strings.Cut(line, token.Colon)
		if !hasSep || strings.TrimSpace(key) != cfgCodex.FrontmatterKeyName {
			continue
		}
		value = strings.Trim(strings.TrimSpace(value), token.Quotes)
		return value == name
	}
	return false
}
