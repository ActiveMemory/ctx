//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

// Meta holds the frontmatter fields extracted from a
// SKILL.md file.
//
// Fields:
//   - Name: Skill name from frontmatter
//   - Description: Skill description from frontmatter
type Meta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
