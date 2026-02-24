//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import "time"

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
	ID   string `json:"id"`
	Slug string `json:"slug,omitempty"`

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

// UserMessages returns only user messages from the session.
//
// Returns:
//   - []Message: Filtered list containing only messages with Role "user"
func (s *Session) UserMessages() []Message {
	var msgs []Message
	for _, m := range s.Messages {
		if m.BelongsToUser() {
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
		if m.BelongsToAssistant() {
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
