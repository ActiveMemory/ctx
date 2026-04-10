//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/claude"
	cfgFmt "github.com/ActiveMemory/ctx/internal/config/format"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/write/guide"
)

// ParseFrontmatter extracts YAML frontmatter from a
// SKILL.md file.
//
// Parameters:
//   - content: Raw SKILL.md content
//
// Returns:
//   - Meta: Parsed name and description (zero if none)
//   - error: Non-nil if YAML parsing fails
func ParseFrontmatter(
	content []byte,
) (Meta, error) {
	text := string(content)
	prefix := token.Separator + token.NewlineLF
	if !strings.HasPrefix(text, prefix) {
		return Meta{}, nil
	}

	offset := len(prefix)
	end := strings.Index(
		text[offset:], token.NewlineLF+token.Separator,
	)
	if end < 0 {
		return Meta{}, nil
	}

	block := []byte(text[offset : offset+end])
	var meta Meta
	if yamlErr := yaml.Unmarshal(block, &meta); yamlErr != nil {
		return Meta{}, yamlErr
	}
	return meta, nil
}

// TruncateDescription returns the first sentence or up to
// maxLen characters.
//
// Parameters:
//   - d: Full description text
//   - maxLen: Maximum character length
//
// Returns:
//   - string: Truncated description
func TruncateDescription(d string, maxLen int) string {
	if idx := strings.Index(d, token.PeriodSpace); idx >= 0 &&
		idx < maxLen {
		return d[:idx+1]
	}
	if len(d) <= maxLen {
		return d
	}
	return d[:maxLen] + token.Ellipsis
}

// List prints all available skills with their descriptions.
//
// Parameters:
//   - cmd: Cobra command for output
//
// Returns:
//   - error: Non-nil if skill listing fails
func List(cmd *cobra.Command) error {
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

		meta, parseErr := ParseFrontmatter(content)
		if parseErr != nil {
			continue
		}

		d := TruncateDescription(
			meta.Description, cfgFmt.TruncateDescription,
		)
		guide.InfoSkillLine(cmd, name, d)
	}
	return nil
}
