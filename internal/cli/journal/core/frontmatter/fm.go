//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package frontmatter

// Obsidian represents the YAML frontmatter for Obsidian vault
// entries. Extends JournalFrontmatter with Obsidian-specific fields.
type Obsidian struct {
	Title        string   `yaml:"title"`
	Date         string   `yaml:"date"`
	Type         string   `yaml:"type,omitempty"`
	Outcome      string   `yaml:"outcome,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Technologies []string `yaml:"technologies,omitempty"`
	KeyFiles     []string `yaml:"key_files,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty"`
	SourceFile   string   `yaml:"source_file,omitempty"`
}
