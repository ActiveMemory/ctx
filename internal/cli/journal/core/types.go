//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// GroupedIndex holds entries aggregated by a string key, sorted by count desc
// then alphabetically. Used by BuildTopicIndex and BuildKeyFileIndex.
type GroupedIndex struct {
	Key     string
	Entries []JournalEntry
	Popular bool
}

// JournalFrontmatter represents YAML frontmatter in enriched journal entries.
type JournalFrontmatter struct {
	Title     string   `yaml:"title"`
	Date      string   `yaml:"date"`
	Time      string   `yaml:"time,omitempty"`
	Project   string   `yaml:"project,omitempty"`
	SessionID string   `yaml:"session_id,omitempty"`
	Model     string   `yaml:"model,omitempty"`
	TokensIn  int      `yaml:"tokens_in,omitempty"`
	TokensOut int      `yaml:"tokens_out,omitempty"`
	Type      string   `yaml:"type"`
	Outcome   string   `yaml:"outcome"`
	Topics    []string `yaml:"topics"`
	KeyFiles  []string `yaml:"key_files"`
	Summary   string   `yaml:"summary,omitempty"`
}

// JournalEntry represents a parsed journal file.
type JournalEntry struct {
	Filename   string
	Title      string
	Date       string
	Time       string
	Project    string
	SessionID  string
	Model      string
	TokensIn   int
	TokensOut  int
	Path       string
	Size       int64
	Suggestive bool
	Topics     []string
	Type       string
	Outcome    string
	KeyFiles   []string
	Summary    string
}

// PopularSplittable is implemented by MOC data types that can be
// split into popular and long-tail groups.
type PopularSplittable interface {
	IsPopular() bool
}

// TopicData holds aggregated data for a single topic.
type TopicData struct {
	Name    string
	Entries []JournalEntry
	Popular bool
}

// IsPopular reports whether this topic is popular.
func (t TopicData) IsPopular() bool { return t.Popular }

// KeyFileData holds aggregated data for a single file path.
type KeyFileData struct {
	Path    string
	Entries []JournalEntry
	Popular bool
}

// IsPopular reports whether this key file is popular.
func (kf KeyFileData) IsPopular() bool { return kf.Popular }

// TypeData holds aggregated data for a session type.
type TypeData struct {
	Name    string
	Entries []JournalEntry
}

// ObsidianFrontmatter represents the YAML frontmatter for Obsidian vault
// entries. Extends JournalFrontmatter with Obsidian-specific fields.
type ObsidianFrontmatter struct {
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
