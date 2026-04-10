//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package parser provides JSONL session file parsing for the recall system.
//
// It parses AI coding assistant session transcripts into structured Go types
// that can be rendered, searched, and analyzed. The package
// uses a tool-agnostic
// Session output type with tool-specific parsers (e.g., ClaudeCode).
package parser

import (
	"encoding/json"
	"time"
)

// ClaudeCode parses Claude Code JSONL session files.
//
// Claude Code stores sessions as JSONL files where each line is a
// self-contained JSON object representing a message. Messages are
// linked via parentUuid and grouped by sessionId.
type ClaudeCode struct{}

// Copilot parses VS Code Copilot Chat JSONL session files.
//
// Copilot Chat stores sessions as JSONL files in VS Code's workspaceStorage
// directory. Each file contains one session. The first line is a full session
// snapshot (kind=0), subsequent lines are incremental patches (kind=1, kind=2).
type Copilot struct{}

// CopilotCLI parses GitHub Copilot CLI session files.
//
// Copilot CLI stores sessions as JSONL files in ~/.copilot/sessions/
// (or $COPILOT_HOME/sessions/). Each file contains one session with
// JSONL-formatted messages similar to Claude Code's format.
type CopilotCLI struct{}

// MarkdownSession parses Markdown session files written by AI agents.
//
// This parser handles the tool-agnostic session format used by non-Claude
// tools (Copilot, Cursor, Aider, etc.) where the AI agent saves session
// summaries as structured Markdown in .context/sessions/.
//
// Expected format:
//
//	# Session: YYYY-MM-DD - Topic
//
//	## What Was Done
//	- ...
//
//	## Decisions
//	- ...
//
//	## Learnings
//	- ...
//
//	## Next Steps
//	- ...
type MarkdownSession struct{}

// Claude Code JSONL raw types.
//
// These types mirror the on-disk JSONL format produced by Claude Code.
// Each line in a Claude Code session file is a self-contained JSON object
// that deserializes into claudeRawMessage.

// claudeRawMessage represents a single JSONL line from a
// Claude Code session.
//
// Fields:
//
// Core:
//   - UUID: Unique message identifier
//   - ParentUUID: Parent message for threading
//   - SessionID: Groups messages into a single session
//   - RequestID: API request correlation identifier
//   - Timestamp: When the message was created
//   - Type: Message role ("user", "assistant", or system)
//   - UserType: Sub-type for user messages
//   - IsSidechain: True if on a sidechain branch
//   - CWD: Working directory at message time
//   - GitBranch: Active git branch at message time
//   - Version: Claude Code version
//   - Slug: URL-friendly session ID (deprecated)
//   - Message: Nested content payload
//
// Envelope (optional, added in newer CC versions):
//   - PlanContent: Full plan-mode document text
//   - IsApiErrorMessage: True for API error responses
//   - SourceToolAssistantUUID: Links tool result to caller
//   - ToolUseResult: CC-level tool error string
//   - Entrypoint: How CC was launched (cli, ide, sdk)
//   - Origin: Message injection source
//   - Error: Error object on failed API responses
//   - ApiError: API-level error object
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

	PlanContent       string `json:"planContent,omitempty"`
	IsApiErrorMessage bool   `json:"isApiErrorMessage,omitempty"`
	ToolUseResult     string `json:"toolUseResult,omitempty"`
	Entrypoint        string `json:"entrypoint,omitempty"`

	SourceToolAssistantUUID string `json:"sourceToolAssistantUUID,omitempty"`

	Origin   json.RawMessage `json:"origin,omitempty"`
	Error    json.RawMessage `json:"error,omitempty"`
	ApiError json.RawMessage `json:"apiError,omitempty"`
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

// section holds a heading and its body content in document order.
type section struct {
	heading string
	body    string
}
