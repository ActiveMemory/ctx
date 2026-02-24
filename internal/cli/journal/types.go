//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package journal

// groupedIndex holds entries aggregated by a string key, sorted by count desc
// then alphabetically. Used by buildTopicIndex and buildKeyFileIndex.
type groupedIndex struct {
	Key     string
	Entries []journalEntry
	Popular bool
}

// journalFrontmatter represents YAML frontmatter in enriched journal entries.
type journalFrontmatter struct {
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

// journalEntry represents a parsed journal file.
type journalEntry struct {
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

// topicData holds aggregated data for a single topic.
type topicData struct {
	Name    string
	Entries []journalEntry
	Popular bool
}

// keyFileData holds aggregated data for a single file path.
type keyFileData struct {
	Path    string
	Entries []journalEntry
	Popular bool
}

// typeData holds aggregated data for a session type.
type typeData struct {
	Name    string
	Entries []journalEntry
}
