//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parser provides JSONL session file parsing for the recall system.
//
// It parses AI coding assistant session transcripts into structured Go types
// that can be rendered, searched, and analyzed. The package uses a tool-agnostic
// Session output type with tool-specific parsers (e.g., ClaudeCodeParser).
package parser

// ToolUse represents a tool invocation by the assistant.
//
// Fields:
//   - ID: Unique identifier for this tool use
//   - Name: Tool name (e.g., "Bash", "Read", "Write", "Grep")
//   - Input: JSON string of input parameters passed to the tool
type ToolUse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Input string `json:"input"`
}

// ToolResult represents the result of a tool invocation.
//
// Fields:
//   - ToolUseID: ID of the ToolUse this result corresponds to
//   - Content: The tool's output content
//   - IsError: True if the tool execution failed
type ToolResult struct {
	ToolUseID string `json:"tool_use_id"`
	Content   string `json:"content"`
	IsError   bool   `json:"is_error,omitempty"`
}

// SessionParser defines the interface for tool-specific session parsers.
//
// Each AI tool (Claude Code, Aider, Cursor) implements this interface
// to parse its specific format into the common Session type.
type SessionParser interface {
	// ParseFile reads a session file and returns all sessions found.
	// A single file may contain multiple sessions (grouped by session ID).
	ParseFile(path string) ([]*Session, error)

	// ParseLine parses a single line from a session file.
	// Returns nil if the line should be skipped (e.g., non-message lines).
	ParseLine(line []byte) (*Message, string, error) // message, sessionID, error

	// Matches returns true if this parser can handle the given file.
	// Implementations may check file extension, peek at content, etc.
	Matches(path string) bool

	// Tool returns the tool identifier (e.g., "claude-code", "aider").
	Tool() string
}
