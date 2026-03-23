//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

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

// BlockResponse is the JSON output for blocked commands.
type BlockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// Stats holds the fields written to the per-session stats JSONL file.
type Stats struct {
	Timestamp  string `json:"ts"`
	Prompt     int    `json:"prompt"`
	Tokens     int    `json:"tokens"`
	Pct        int    `json:"pct"`
	WindowSize int    `json:"window"`
	Model      string `json:"model,omitempty"`
	Event      string `json:"event"`
}

// TokenInfo holds token usage and model information extracted from a
// session's JSONL file.
type TokenInfo struct {
	Tokens int    // Total input tokens (input + cache_creation + cache_read)
	Model  string // Model ID from the last assistant message, or ""
}
