//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/ActiveMemory/ctx/internal/claude"
)

// skillMeta holds the frontmatter fields we care about.
type skillMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// parseSkillFrontmatter extracts YAML frontmatter from a SKILL.md file.
//
// Returns:
//   - skillMeta: Parsed name and description (zero value if no frontmatter)
//   - error: Non-nil if YAML parsing fails
func parseSkillFrontmatter(content []byte) (skillMeta, error) {
	const sep = "---"

	text := string(content)
	if !strings.HasPrefix(text, sep+"\n") {
		return skillMeta{}, nil
	}

	end := strings.Index(text[4:], "\n"+sep)
	if end < 0 {
		return skillMeta{}, nil
	}

	block := []byte(text[4 : 4+end])
	var meta skillMeta
	if yamlErr := yaml.Unmarshal(block, &meta); yamlErr != nil {
		return skillMeta{}, yamlErr
	}
	return meta, nil
}

// truncateDescription returns the first sentence or up to maxLen characters.
func truncateDescription(desc string, maxLen int) string {
	// Try sentence boundary first.
	if idx := strings.Index(desc, ". "); idx >= 0 && idx < maxLen {
		return desc[:idx+1]
	}
	if len(desc) <= maxLen {
		return desc
	}
	return desc[:maxLen] + "..."
}

// listSkills prints all available skills with their descriptions.
func listSkills(cmd *cobra.Command) error {
	names, skillsErr := claude.Skills()
	if skillsErr != nil {
		return skillsErr
	}

	cmd.Println("Available Skills:")
	cmd.Println()

	var buf bytes.Buffer
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
		buf.Reset()
		fmt.Fprintf(&buf, "  /%-22s %s", name, desc)
		cmd.Println(buf.String())
	}
	return nil
}
