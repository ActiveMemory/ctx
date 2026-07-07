//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package state

// State is the top-level state file structure.
//
// Fields:
//   - Version: Schema version for forward compatibility
//   - Entries: Per-file processing state keyed by filename
//   - Sessions: Per-session source tracking keyed by session ID
//     (schema v2). Records the transcript each session was last
//     rendered from so growth-aware import can detect when the
//     source has changed. Absent in v1 files; initialised empty on
//     load.
type State struct {
	Version  int               `json:"version"`
	Entries  map[string]File   `json:"entries"`
	Sessions map[string]Source `json:"sessions,omitempty"`
}

// File tracks processing stages for a single journal entry.
// Values are date strings (YYYY-MM-DD) indicating when the stage completed.
//
// Fields:
//   - Exported: Date the session was imported to markdown
//   - Enriched: Date frontmatter was added
//   - Normalized: Date formatting was cleaned up
//   - FencesVerified: Date fence balance was checked
//   - Locked: Date the entry was locked from regeneration
//   - RenderHash: Hash of the last ctx-authored write of the file
//     (schema v2). Lets growth-aware import prove a file is
//     unedited before splicing fresh transcript into it; a mismatch
//     means a human edited the file and it must not be clobbered.
//     Empty in v1 files and for entries never written under v2.
type File struct {
	Exported       string `json:"exported,omitempty"`
	Enriched       string `json:"enriched,omitempty"`
	Normalized     string `json:"normalized,omitempty"`
	FencesVerified string `json:"fences_verified,omitempty"`
	Locked         string `json:"locked,omitempty"`
	RenderHash     string `json:"render_hash,omitempty"`
}

// Source records the transcript a session was last rendered from, so
// growth-aware import can detect when the source has changed.
//
// Claude Code JSONL transcripts are append-only, so a change in mtime
// or size is sufficient to detect growth — no content hashing of the
// source is needed. Keyed by session ID (not source path) because one
// session can span multiple transcript files: a resume copies prior
// history into a new file, and switching to that larger copy also
// counts as growth.
//
// Fields:
//   - SourceFile: Absolute path to the transcript last rendered from
//   - SourceMtime: Unix mtime (seconds) of that transcript at render time
//   - SourceSize: Byte size of that transcript at render time
type Source struct {
	SourceFile  string `json:"source_file"`
	SourceMtime int64  `json:"source_mtime"`
	SourceSize  int64  `json:"source_size"`
}
