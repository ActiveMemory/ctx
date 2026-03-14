//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/claude"
	"github.com/ActiveMemory/ctx/internal/write"
)

// parseSkillFrontmatter extracts YAML frontmatter from a SKILL.md file.
//
// Parameters:
//   - content: Raw SKILL.md content
//
// Returns:
//   - skillMeta: Parsed name and description (zero value if no frontmatter)
//   - error: Non-nil if YAML parsing fails
func parseSkillFrontmatter(content []byte) (skillMeta, error) {
	text := string(content)
	prefix := token.Separator + token.NewlineLF
	if !strings.HasPrefix(text, prefix) {
		return skillMeta{}, nil
	}

	offset := len(prefix)
	end := strings.Index(text[offset:], token.NewlineLF+token.Separator)
	if end < 0 {
		return skillMeta{}, nil
	}

	block := []byte(text[offset : offset+end])
	var meta skillMeta
	if yamlErr := yaml.Unmarshal(block, &meta); yamlErr != nil {
		return skillMeta{}, yamlErr
	}
	return meta, nil
}

// truncateDescription returns the first sentence or up to maxLen characters.
//
// Parameters:
//   - desc: Full description text
//   - maxLen: Maximum character length
//
// Returns:
//   - string: Truncated description
func truncateDescription(desc string, maxLen int) string {
	if idx := strings.Index(desc, ". "); idx >= 0 && idx < maxLen {
		return desc[:idx+1]
	}
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + token.Ellipsis
}

// listSkills prints all available skills with their descriptions.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if skill listing fails
func listSkills(cmd *cobra.Command) error {
	names, skillsErr := claude.Skills()
	if skillsErr != nil {
		return skillsErr
	}

	write.InfoSkillsHeader(cmd)

	for _, name := range names {
		content, readErr := claude.SkillContent(name)
		if readErr != nil {
			continue
		}

		meta, parseErr := parseSkillFrontmatter(content)
		if parseErr != nil {
			continue
		}

		desc := truncateDescription(meta.Description, 70)
		write.InfoSkillLine(cmd, name, desc)
	}
	return nil
}
