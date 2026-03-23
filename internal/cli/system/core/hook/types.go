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

// ResourceJSONAlert is a single resource alert for JSON output.
type ResourceJSONAlert struct {
	Severity string `json:"severity"`
	Resource string `json:"resource"`
	Message  string `json:"message"`
}

// ResourceJSONOutput is the top-level JSON output for system resources.
type ResourceJSONOutput struct {
	Memory struct {
		TotalBytes uint64 `json:"total_bytes"`
		UsedBytes  uint64 `json:"used_bytes"`
		Percent    int    `json:"percent"`
		Supported  bool   `json:"supported"`
	} `json:"memory"`
	Swap struct {
		TotalBytes uint64 `json:"total_bytes"`
		UsedBytes  uint64 `json:"used_bytes"`
		Percent    int    `json:"percent"`
		Supported  bool   `json:"supported"`
	} `json:"swap"`
	Disk struct {
		TotalBytes uint64 `json:"total_bytes"`
		UsedBytes  uint64 `json:"used_bytes"`
		Percent    int    `json:"percent"`
		Path       string `json:"path"`
		Supported  bool   `json:"supported"`
	} `json:"disk"`
	Load struct {
		Load1     float64 `json:"load1"`
		Load5     float64 `json:"load5"`
		Load15    float64 `json:"load15"`
		NumCPU    int     `json:"num_cpu"`
		Ratio     float64 `json:"ratio"`
		Supported bool    `json:"supported"`
	} `json:"load"`
	Alerts      []ResourceJSONAlert `json:"alerts"`
	MaxSeverity string              `json:"max_severity"`
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
