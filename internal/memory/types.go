//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import "time"

// EntryKind identifies how an entry was delimited in MEMORY.md.
type EntryKind int

const (
	// EntryHeader is a Markdown heading (## or ###).
	EntryHeader EntryKind = iota
	// EntryParagraph is a blank-line-separated paragraph.
	EntryParagraph
	// EntryList is one or more consecutive list items.
	EntryList
)

// Entry is a discrete block parsed from MEMORY.md.
type Entry struct {
	Text      string    // Raw text of the entry (trimmed)
	StartLine int       // 1-based line number where the entry begins
	Kind      EntryKind // How the entry was delimited
}

// Classification is the result of heuristic entry classification.
type Classification struct {
	Target   string   // config.Entry* constant or "skip"
	Keywords []string // Keywords that triggered the match
}

// PublishResult holds what was selected for publishing.
//
// Fields:
//   - Tasks: Task entries selected for MEMORY.md
//   - Decisions: Decision entries selected
//   - Conventions: Convention entries selected
//   - Learnings: Learning entries selected
//   - TotalLines: Total lines across all selections
type PublishResult struct {
	Tasks       []string
	Decisions   []string
	Conventions []string
	Learnings   []string
	TotalLines  int
}

// State tracks memory bridge sync timestamps and import/publish progress.
//
// Fields:
//   - LastSync: When mirror was last updated
//   - LastImport: When entries were last imported from MEMORY.md
//   - LastPublish: When context was last published to MEMORY.md
//   - ImportedHashes: Content hashes of already-imported entries
type State struct {
	LastSync       *time.Time `json:"last_sync"`
	LastImport     *time.Time `json:"last_import"`
	LastPublish    *time.Time `json:"last_publish"`
	ImportedHashes []string   `json:"imported_hashes"`
}

// SyncResult holds the outcome of a Sync operation.
//
// Fields:
//   - SourcePath: Path to the source MEMORY.md
//   - MirrorPath: Path to the mirror copy
//   - ArchivedTo: Archive path (empty if no prior mirror)
//   - SourceLines: Line count of the source file
//   - MirrorLines: Line count of the previous mirror (0 if first sync)
type SyncResult struct {
	SourcePath  string
	MirrorPath  string
	ArchivedTo  string // empty if no prior mirror existed
	SourceLines int
	MirrorLines int // lines in the previous mirror (0 if first sync)
}
