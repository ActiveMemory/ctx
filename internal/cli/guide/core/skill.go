//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/claude"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/write/guide"
)

// ParseSkillFrontmatter extracts YAML frontmatter from a
// SKILL.md file.
//
// Parameters:
//   - content: Raw SKILL.md content
//
// Returns:
//   - SkillMeta: Parsed name and description (zero if none)
//   - error: Non-nil if YAML parsing fails
func ParseSkillFrontmatter(
	content []byte,
) (SkillMeta, error) {
	text := string(content)
	prefix := token.Separator + token.NewlineLF
	if !strings.HasPrefix(text, prefix) {
		return SkillMeta{}, nil
	}

	offset := len(prefix)
	end := strings.Index(
		text[offset:], token.NewlineLF+token.Separator,
	)
	if end < 0 {
		return SkillMeta{}, nil
	}

	block := []byte(text[offset : offset+end])
	var meta SkillMeta
	if yamlErr := yaml.Unmarshal(block, &meta); yamlErr != nil {
		return SkillMeta{}, yamlErr
	}
	return meta, nil
}

// TruncateDescription returns the first sentence or up to
// maxLen characters.
//
// Parameters:
//   - desc: Full description text
//   - maxLen: Maximum character length
//
// Returns:
//   - string: Truncated description
func TruncateDescription(desc string, maxLen int) string {
	if idx := strings.Index(desc, ". "); idx >= 0 &&
		idx < maxLen {
		return desc[:idx+1]
	}
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + token.Ellipsis
}

// ListSkills prints all available skills with their
// descriptions.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if skill listing fails
func ListSkills(cmd *cobra.Command) error {
	names, skillsErr := claude.SkillList()
	if skillsErr != nil {
		return skillsErr
	}

	guide.InfoSkillsHeader(cmd)

	for _, name := range names {
		content, readErr := claude.SkillContent(name)
		if readErr != nil {
			continue
		}

		meta, parseErr := ParseSkillFrontmatter(content)
		if parseErr != nil {
			continue
		}

		desc := TruncateDescription(
			meta.Description, cfgFmt.TruncateDescription,
		)
		guide.InfoSkillLine(cmd, name, desc)
	}
	return nil
}
