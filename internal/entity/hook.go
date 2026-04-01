//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// HookInput represents the JSON payload that Claude Code sends to hook
// commands via stdin.
//
// Fields:
//   - SessionID: Claude Code session identifier
//   - ToolInput: Tool-specific fields from the invocation
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
//
// Fields:
//   - Decision: "block" or "allow"
//   - Reason: Human-readable explanation
type BlockResponse struct {
	Decision string `json:"decision"`
	Reason   string `json:"reason"`
}

// Stats holds the fields written to the per-session stats JSONL file.
//
// Fields:
//   - Timestamp: ISO 8601 timestamp
//   - Prompt: Prompt counter within the session
//   - Tokens: Total token count at this point
//   - Pct: Percentage of context window used
//   - WindowSize: Context window size in tokens
//   - Model: Model ID (omitted if unknown)
//   - Event: Event type that triggered this entry
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
//
// Fields:
//   - Tokens: Total input tokens (input + cache_creation + cache_read)
//   - Model: Model ID from the last assistant message, or ""
type TokenInfo struct {
	Tokens int    // Total input tokens (input + cache_creation + cache_read)
	Model  string // Model ID from the last assistant message, or ""
}
