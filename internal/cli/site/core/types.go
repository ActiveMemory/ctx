//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// BlogPost holds parsed metadata from a single blog post.
type BlogPost struct {
	Filename string
	Title    string
	Date     string
	Author   string
	Topics   []string
	Summary  string
}

// FeedReport tracks what happened during feed generation.
type FeedReport struct {
	Included int
	Skipped  []string
	Warnings []string
}

// BlogFrontmatter maps the YAML fields we care about.
type BlogFrontmatter struct {
	Title                string   `yaml:"title"`
	Date                 string   `yaml:"date"`
	Author               string   `yaml:"author"`
	Topics               []string `yaml:"topics"`
	ReviewedAndFinalized *bool    `yaml:"reviewed_and_finalized"`
}

// PostStatus indicates the outcome of parsing a single blog post.
type PostStatus int

const (
	PostIncluded PostStatus = iota
	PostSkipped
	PostWarn
)
