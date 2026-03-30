//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// JournalFrontmatter represents YAML frontmatter in enriched journal entries.
//
// Fields:
//   - Title: Session title
//   - Date: Date string (YYYY-MM-DD)
//   - Time: Time string (HH:MM, optional)
//   - Project: Project name
//   - SessionID: Claude Code session UUID
//   - Model: Model ID used in session
//   - TokensIn: Input tokens consumed
//   - TokensOut: Output tokens generated
//   - Type: Session type (feature, debug, refactor, etc.)
//   - Outcome: Session outcome (completed, partial, etc.)
//   - Topics: Topic tags for indexing
//   - KeyFiles: Files referenced in the session
//   - Summary: One-paragraph summary
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
//
// Fields:
//   - Filename: Journal markdown filename
//   - Title: Session title from frontmatter
//   - Date: Date string from frontmatter
//   - Time: Time string from frontmatter
//   - Project: Project name
//   - SessionID: Claude Code session UUID
//   - Model: Model ID
//   - TokensIn: Input tokens
//   - TokensOut: Output tokens
//   - Path: Full file path
//   - Size: File size in bytes
//   - Suggestive: Whether the title was auto-generated
//   - Topics: Topic tags
//   - Type: Session type
//   - Outcome: Session outcome
//   - KeyFiles: Referenced file paths
//   - Summary: One-paragraph summary
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

// GroupedIndex holds entries aggregated by a string key, sorted by count desc
// then alphabetically. Used by BuildTopicIndex and BuildKeyFileIndex.
//
// Fields:
//   - Key: Grouping key (topic name or file path)
//   - Entries: Journal entries in this group
//   - Popular: Whether this group exceeds the popularity threshold
type GroupedIndex struct {
	Key     string
	Entries []JournalEntry
	Popular bool
}

// PopularSplittable is implemented by MOC data types that can be
// split into popular and long-tail groups.
type PopularSplittable interface {
	IsPopular() bool
}

// TopicData holds aggregated data for a single topic.
//
// Fields:
//   - Name: Topic name
//   - Entries: Journal entries tagged with this topic
//   - Popular: Whether this topic is in the popular set
type TopicData struct {
	Name    string
	Entries []JournalEntry
	Popular bool
}

// IsPopular reports whether this topic is popular.
func (t TopicData) IsPopular() bool { return t.Popular }

// KeyFileData holds aggregated data for a single file path.
//
// Fields:
//   - Path: File path
//   - Entries: Journal entries referencing this file
//   - Popular: Whether this file is in the popular set
type KeyFileData struct {
	Path    string
	Entries []JournalEntry
	Popular bool
}

// IsPopular reports whether this key file is popular.
func (kf KeyFileData) IsPopular() bool { return kf.Popular }

// TypeData holds aggregated data for a session type.
//
// Fields:
//   - Name: Session type name (feature, debug, refactor, etc.)
//   - Entries: Journal entries of this type
type TypeData struct {
	Name    string
	Entries []JournalEntry
}
