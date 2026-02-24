//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import "time"

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

// BelongsToUser returns true if this is a user message.
//
// Returns:
//   - bool: True if Role is "user"
func (m *Message) BelongsToUser() bool {
	return m.Role == "user"
}

// BelongsToAssistant returns true if this is an assistant message.
//
// Returns:
//   - bool: True if Role is "assistant"
func (m *Message) BelongsToAssistant() bool {
	return m.Role == "assistant"
}

// UsesTools returns true if this message contains tool invocations.
//
// Returns:
//   - bool: True if ToolUses slice is non-empty
func (m *Message) UsesTools() bool {
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
