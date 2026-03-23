//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

// ArchiveEntry describes a directory or file to include in a backup archive.
type ArchiveEntry struct {
	// SourcePath is the absolute path to the directory or file.
	SourcePath string
	// Prefix is the path prefix inside the tar archive.
	Prefix string
	// ExcludeDir is a directory name to skip (e.g. "journal-site").
	ExcludeDir string
	// Optional means a missing source is not an error.
	Optional bool
}

// BackupResult holds the outcome of a single archive creation.
type BackupResult struct {
	Scope   string `json:"scope"`
	Archive string `json:"archive"`
	Size    int64  `json:"size"`
	SMBDest string `json:"smb_dest,omitempty"`
}

// FileTokenEntry tracks per-file token counts during context injection.
type FileTokenEntry struct {
	Name   string
	Tokens int
}

// MessageListEntry holds the data for a single row in the message list output.
type MessageListEntry struct {
	Hook         string   `json:"hook"`
	Variant      string   `json:"variant"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	TemplateVars []string `json:"template_vars"`
	HasOverride  bool     `json:"has_override"`
}

// StaleEntry describes a file that has not been modified within the
// configured freshness window.
type StaleEntry struct {
	// Path is the relative file path.
	Path string
	// Desc is the human-readable file description.
	Desc string
	// ReviewURL is the optional URL for reviewing the file against upstream.
	ReviewURL string
	// Days is the number of days since last modification.
	Days int
}
