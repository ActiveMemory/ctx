//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

// Skill represents a parsed skill manifest (SKILL.md) with YAML
// frontmatter and markdown instruction body.
//
// Fields:
//   - Name: Unique identifier from frontmatter
//   - Description: Short summary from frontmatter
//   - Body: Markdown instruction content after frontmatter
//   - Dir: Directory path containing the SKILL.md file
type Skill struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Body        string `yaml:"-"`
	Dir         string `yaml:"-"`
}
