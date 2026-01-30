//   /    Context:                     https://ctx.ist
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

import (
	"time"
)

// Session represents a reconstructed conversation session.
//
// This is the tool-agnostic output type that all parsers produce.
// It contains common fields that make sense across different AI tools.
//
// Fields:
//
// Identity:
//   - ID: Unique session identifier
//   - Slug: URL-friendly session identifier
//
// Source:
//   - Tool: Source tool ("claude-code", "aider", "cursor", etc.)
//   - SourceFile: Original file path
//
// Context:
//   - CWD: Working directory when session started
//   - Project: Project name (derived from last component of CWD)
//   - GitBranch: Git branch name if available
//
// Timing:
//   - StartTime: When the session started
//   - EndTime: When the session ended
//   - Duration: Total session duration
//
// Messages:
//   - Messages: All messages in the session
//   - TurnCount: Count of user messages
//
// Token Statistics:
//   - TotalTokensIn: Input tokens used (if available)
//   - TotalTokensOut: Output tokens used (if available)
//   - TotalTokens: Total tokens used (if available)
//
// Derived:
//   - HasErrors: True if any tool errors occurred
//   - FirstUserMsg: Preview text of first user message (truncated)
//   - Model: Primary model used in the session
type Session struct {
	ID         string `json:"id"`
	Slug       string `json:"slug,omitempty"`

	Tool       string `json:"tool"`
	SourceFile string `json:"source_file"`

	CWD       string `json:"cwd,omitempty"`
	Project   string `json:"project,omitempty"`
	GitBranch string `json:"git_branch,omitempty"`

	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`

	Messages  []Message `json:"messages"`
	TurnCount int       `json:"turn_count"`

	TotalTokensIn  int `json:"total_tokens_in,omitempty"`
	TotalTokensOut int `json:"total_tokens_out,omitempty"`
	TotalTokens    int `json:"total_tokens,omitempty"`

	HasErrors    bool   `json:"has_errors,omitempty"`
	FirstUserMsg string `json:"first_user_msg,omitempty"`
	Model        string `json:"model,omitempty"`
}

// Message represents a single message in a session.
//
// This is tool-agnostic - all parsers normalize to this format.
//
// Fields:
//
// Identity:
//   - ID: Unique message identifier
//   - Timestamp: When the message was created
//   - Role: Message role ("user" or "assistant")
//
// Content:
//   - Text: Main text content
//   - Thinking: Reasoning content (if available)
//   - ToolUses: Tool invocations in this message
//   - ToolResults: Results from tool invocations
//
// Token Usage:
//   - TokensIn: Input tokens for this message (if available)
//   - TokensOut: Output tokens for this message (if available)
type Message struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Role      string    `json:"role"`

	Text        string       `json:"text,omitempty"`
	Thinking    string       `json:"thinking,omitempty"`
	ToolUses    []ToolUse    `json:"tool_uses,omitempty"`
	ToolResults []ToolResult `json:"tool_results,omitempty"`

	TokensIn  int `json:"tokens_in,omitempty"`
	TokensOut int `json:"tokens_out,omitempty"`
}

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

	// CanParse returns true if this parser can handle the given file.
	// Implementations may check file extension, peek at content, etc.
	CanParse(path string) bool

	// Tool returns the tool identifier (e.g., "claude-code", "aider").
	Tool() string
}

// IsUser returns true if this is a user message.
//
// Returns:
//   - bool: True if Role is "user"
func (m *Message) IsUser() bool {
	return m.Role == "user"
}

// IsAssistant returns true if this is an assistant message.
//
// Returns:
//   - bool: True if Role is "assistant"
func (m *Message) IsAssistant() bool {
	return m.Role == "assistant"
}

// HasToolUses returns true if this message contains tool invocations.
//
// Returns:
//   - bool: True if ToolUses slice is non-empty
func (m *Message) HasToolUses() bool {
	return len(m.ToolUses) > 0
}

// Preview returns a truncated preview of the message text.
//
// Parameters:
//   - maxLen: Maximum length before truncation (adds "..." if exceeded)
//
// Returns:
//   - string: The text, truncated with "..." suffix if longer than maxLen
func (m *Message) Preview(maxLen int) string {
	if len(m.Text) <= maxLen {
		return m.Text
	}
	return m.Text[:maxLen] + "..."
}

// UserMessages returns only user messages from the session.
//
// Returns:
//   - []Message: Filtered list containing only messages with Role "user"
func (s *Session) UserMessages() []Message {
	var msgs []Message
	for _, m := range s.Messages {
		if m.IsUser() {
			msgs = append(msgs, m)
		}
	}
	return msgs
}

// AssistantMessages returns only assistant messages from the session.
//
// Returns:
//   - []Message: Filtered list containing only messages with Role "assistant"
func (s *Session) AssistantMessages() []Message {
	var msgs []Message
	for _, m := range s.Messages {
		if m.IsAssistant() {
			msgs = append(msgs, m)
		}
	}
	return msgs
}

// AllToolUses returns all tool uses across all messages.
//
// Returns:
//   - []ToolUse: Aggregated list of all tool invocations in the session
func (s *Session) AllToolUses() []ToolUse {
	var tools []ToolUse
	for _, m := range s.Messages {
		tools = append(tools, m.ToolUses...)
	}
	return tools
}
