//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package feed

// blogPost holds parsed metadata from a single blog post.
type blogPost struct {
	filename string
	title    string
	date     string
	author   string
	topics   []string
	summary  string
}

// feedReport tracks what happened during feed generation.
type feedReport struct {
	included int
	skipped  []string // "filename — reason"
	warnings []string // "filename — reason"
}

// blogFrontmatter maps the YAML fields we care about.
type blogFrontmatter struct {
	Title                string   `yaml:"title"`
	Date                 string   `yaml:"date"`
	Author               string   `yaml:"author"`
	Topics               []string `yaml:"topics"`
	ReviewedAndFinalized *bool    `yaml:"reviewed_and_finalized"`
}
