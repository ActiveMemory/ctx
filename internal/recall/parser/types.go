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

import (
	"encoding/json"
	"time"
)

// Claude Code JSONL raw types.
//
// These types mirror the on-disk JSONL format produced by Claude Code.
// Each line in a Claude Code session file is a self-contained JSON object
// that deserializes into claudeRawMessage.

// claudeRawMessage represents a single JSONL line from a Claude Code session.
//
// Fields:
//   - UUID: Unique message identifier
//   - ParentUUID: Parent message for threading (nil for root messages)
//   - SessionID: Groups messages into a single session
//   - RequestID: API request correlation identifier
//   - Timestamp: When the message was created
//   - Type: Message role ("user", "assistant", or system types)
//   - UserType: Sub-type for user messages
//   - IsSidechain: True if the message is on a sidechain branch
//   - CWD: Working directory at message time
//   - GitBranch: Active git branch at message time
//   - Version: Claude Code version that created the message
//   - Slug: URL-friendly session identifier (removed in newer versions)
//   - Message: Nested content payload
type claudeRawMessage struct {
	UUID        string           `json:"uuid"`
	ParentUUID  *string          `json:"parentUuid"`
	SessionID   string           `json:"sessionId"`
	RequestID   string           `json:"requestId,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
	Type        string           `json:"type"`
	UserType    string           `json:"userType,omitempty"`
	IsSidechain bool             `json:"isSidechain,omitempty"`
	CWD         string           `json:"cwd"`
	GitBranch   string           `json:"gitBranch,omitempty"`
	Version     string           `json:"version"`
	Slug        string           `json:"slug"`
	Message     claudeRawContent `json:"message"`
}

// claudeRawContent is the nested content envelope inside a claudeRawMessage.
//
// Fields:
//   - ID: Content block identifier
//   - Type: Content type discriminator
//   - Model: AI model used for this response
//   - Role: Message role ("user" or "assistant")
//   - Content: Raw JSON that may be a string or []claudeRawBlock
//   - StopReason: Why the model stopped generating
//   - StopSequence: Stop sequence that was hit, if any
//   - Usage: Token usage statistics for this message
type claudeRawContent struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`
	Model        string          `json:"model,omitempty"`
	Role         string          `json:"role"`
	Content      json.RawMessage `json:"content"`
	StopReason   *string         `json:"stop_reason,omitempty"`
	StopSequence *string         `json:"stop_sequence,omitempty"`
	Usage        *claudeRawUsage `json:"usage,omitempty"`
}

// claudeRawBlock represents a single content block in a Claude response.
//
// The Type field discriminates between text, thinking, tool_use, and
// tool_result blocks. Only fields relevant to the block type are populated.
//
// Fields:
//   - Type: Block type ("text", "thinking", "tool_use", "tool_result")
//   - Text: Text content (for text blocks)
//   - Thinking: Reasoning content (for thinking blocks)
//   - Signature: Cryptographic signature (for thinking blocks)
//   - ID: Block identifier (for tool_use blocks)
//   - Name: Tool name (for tool_use blocks)
//   - Input: Raw JSON tool parameters (for tool_use blocks)
//   - ToolUseID: References the tool_use block (for tool_result blocks)
//   - Content: Raw JSON tool output (for tool_result blocks)
//   - IsError: True if tool execution failed (for tool_result blocks)
type claudeRawBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	Thinking  string          `json:"thinking,omitempty"`
	Signature string          `json:"signature,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
	IsError   bool            `json:"is_error,omitempty"`
}

// claudeRawUsage contains token usage statistics from the Claude API.
//
// Fields:
//   - InputTokens: Number of input tokens consumed
//   - OutputTokens: Number of output tokens generated
//   - CacheCreationInputTokens: Tokens used to create prompt cache
//   - CacheReadInputTokens: Tokens read from prompt cache
type claudeRawUsage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
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

// section holds a heading and its body content in document order.
type section struct {
	heading string
	body    string
}
