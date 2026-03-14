//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// BelongsToUser returns true if this is a user message.
//
// Returns:
//   - bool: True if Role is "user"
func (m *Message) BelongsToUser() bool {
	return m.Role == claude.RoleUser
}

// BelongsToAssistant returns true if this is an assistant message.
//
// Returns:
//   - bool: True if Role is "assistant"
func (m *Message) BelongsToAssistant() bool {
	return m.Role == claude.RoleAssistant
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
	return m.Text[:maxLen] + token.Ellipsis
}
