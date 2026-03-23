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

// BlockResponse is the JSON output for blocked commands.
type BlockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// HookInput represents the JSON payload that Claude Code sends to hook
// commands via stdin.
type HookInput struct {
	SessionID string    `json:"session_id"`
	ToolInput ToolInput `json:"tool_input"`
}

// ToolInput contains the tool-specific fields from a Claude Code hook
// invocation. For Bash hooks, Command holds the shell command.
type ToolInput struct {
	Command string `json:"command"`
}

// HookResponse is the JSON output format for Claude Code hooks.
// Using structured JSON ensures the agent processes the output as a directive
// rather than treating it as ignorable plain text.
type HookResponse struct {
	HookSpecificOutput *HookSpecificOutput `json:"hookSpecificOutput,omitempty"`
}

// HookSpecificOutput carries event-specific fields inside a HookResponse.
type HookSpecificOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext,omitempty"`
}

// FileTokenEntry tracks per-file token counts during context injection.
type FileTokenEntry struct {
	Name   string
	Tokens int
}

// StatsEntry is a SessionStats with the source file for display.
type StatsEntry struct {
	SessionStats
	Session string `json:"session"`
}

// SessionTokenInfo holds token usage and model information extracted from a
// session's JSONL file.
type SessionTokenInfo struct {
	Tokens int    // Total input tokens (input + cache_creation + cache_read)
	Model  string // Model ID from the last assistant message, or ""
}

// usageData represents the minimal usage fields from a Claude Code JSONL
// assistant message. Only the fields needed for token counting are included.
type usageData struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// jsonlMessage represents the minimal structure of a Claude Code JSONL line
// needed to extract usage and model data from assistant messages.
type jsonlMessage struct {
	Type    string `json:"type"`
	Message struct {
		Role  string    `json:"role"`
		Model string    `json:"model"`
		Usage usageData `json:"usage"`
	} `json:"message"`
}

// PersistenceState holds the counter state for persistence nudging.
type PersistenceState struct {
	Count     int
	LastNudge int
	LastMtime int64
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

// MapTrackingInfo holds the minimal fields needed from map-tracking.json.
type MapTrackingInfo struct {
	OptedOut bool   `json:"opted_out"`
	LastRun  string `json:"last_run"`
}

// KnowledgeFinding describes a single knowledge file that exceeds its
// configured threshold.
type KnowledgeFinding struct {
	// File is the context filename (e.g., DECISIONS.md).
	File string
	// Count is the actual entry or line count.
	Count int
	// Threshold is the configured maximum.
	Threshold int
	// Unit is the measurement unit ("entries" or "lines").
	Unit string
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
