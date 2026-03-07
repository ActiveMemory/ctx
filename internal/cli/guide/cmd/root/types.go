//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

// skillMeta holds the frontmatter fields extracted from a SKILL.md file.
type skillMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
