//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

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

// Codex parses OpenAI Codex rollout transcripts.
//
// Codex stores one session per JSONL file under
// `$CODEX_HOME/sessions/YYYY/MM/DD/rollout-<ISO>-<uuid>.jsonl`.
// The first line is a `session_meta` envelope carrying the session
// identity, cwd, and git state; the remaining lines are
// `response_item` (messages, tool calls, tool outputs),
// `event_msg` (UI events, token counts), `turn_context`
// (per-turn model), and bookkeeping lines (`world_state`,
// `compacted`) that carry no conversation content.
type Codex struct{}

// Codex rollout JSONL raw types.
//
// These types mirror the on-disk rollout format produced by Codex
// CLI. Every line deserializes into codexRawLine; the payload shape
// depends on the line type.

// codexRawLine is the envelope shared by every rollout line.
//
// Fields:
//   - Timestamp: when the line was written (RFC 3339)
//   - Type: line discriminator (session_meta, response_item, ...)
//   - Payload: type-specific body, decoded lazily
type codexRawLine struct {
	Timestamp time.Time       `json:"timestamp"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
}

// codexRawSessionMeta is the payload of a session_meta line.
//
// Fields:
//   - ID: session identifier (newer releases)
//   - SessionID: session identifier (older releases; fallback)
//   - Timestamp: session start time
//   - CWD: working directory the session was started in
//   - Originator: launching surface (codex_cli_rs, codex_exec, ...)
//   - CLIVersion: Codex CLI version that wrote the rollout
//   - Git: repository state at session start, if any
type codexRawSessionMeta struct {
	ID         string       `json:"id"`
	SessionID  string       `json:"session_id"`
	Timestamp  time.Time    `json:"timestamp"`
	CWD        string       `json:"cwd"`
	Originator string       `json:"originator"`
	CLIVersion string       `json:"cli_version"`
	Git        *codexRawGit `json:"git,omitempty"`
}

// codexRawGit is the git block inside session_meta.
//
// Fields:
//   - Branch: checked-out branch
//   - CommitHash: HEAD commit
//   - RepositoryURL: origin remote URL
type codexRawGit struct {
	Branch        string `json:"branch"`
	CommitHash    string `json:"commit_hash"`
	RepositoryURL string `json:"repository_url"`
}

// codexRawTurnContext is the payload of a turn_context line.
//
// Fields:
//   - Model: model slug used for the turn
//   - CWD: working directory for the turn
type codexRawTurnContext struct {
	Model string `json:"model"`
	CWD   string `json:"cwd"`
}

// codexRawItem is the payload of a response_item line.
//
// The Type field discriminates between message, function_call,
// function_call_output, custom_tool_call, custom_tool_call_output,
// local_shell_call, and reasoning items. Only the fields relevant
// to the item type are populated.
//
// Fields:
//   - Type: item discriminator
//   - ID: item identifier
//   - Role: message role (user, assistant, developer)
//   - Content: message content parts (for message items)
//   - Name: tool name (for function_call / custom_tool_call)
//   - Arguments: JSON-encoded arguments (for function_call)
//   - Input: free-form input (for custom_tool_call)
//   - Action: shell action object (for local_shell_call)
//   - CallID: correlates a call with its output
//   - Output: tool output; a string or an array of content parts
type codexRawItem struct {
	Type      string          `json:"type"`
	ID        string          `json:"id,omitempty"`
	Role      string          `json:"role,omitempty"`
	Content   []codexRawPart  `json:"content,omitempty"`
	Name      string          `json:"name,omitempty"`
	Arguments string          `json:"arguments,omitempty"`
	Input     string          `json:"input,omitempty"`
	Action    json.RawMessage `json:"action,omitempty"`
	CallID    string          `json:"call_id,omitempty"`
	Output    json.RawMessage `json:"output,omitempty"`
}

// codexRawPart is a single content part inside a message or a
// tool output.
//
// Fields:
//   - Type: part discriminator (input_text, output_text)
//   - Text: the part text
type codexRawPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// codexRawEvent is the payload of an event_msg line.
//
// Only token_count events are consumed; the Info block is nil for
// every other event type.
//
// Fields:
//   - Type: event discriminator
//   - Info: token usage block (token_count only)
type codexRawEvent struct {
	Type string             `json:"type"`
	Info *codexRawTokenInfo `json:"info,omitempty"`
	Item *codexRawEventItem `json:"item,omitempty"`
}

// codexRawEventItem is the completed item inside an item_completed
// event; only the error-signal fields are decoded.
//
// Fields:
//   - Type: item discriminator (e.g. CommandExecution)
//   - ExitCode: command exit status, when the item is a command
type codexRawEventItem struct {
	Type     string `json:"type"`
	ExitCode *int   `json:"exit_code,omitempty"`
}

// codexRawTokenInfo is the info block of a token_count event.
//
// Fields:
//   - TotalTokenUsage: cumulative usage for the whole session
type codexRawTokenInfo struct {
	TotalTokenUsage codexRawTokenUsage `json:"total_token_usage"`
}

// codexRawTokenUsage is a token usage record.
//
// Fields:
//   - InputTokens: prompt tokens (includes cached)
//   - CachedInputTokens: prompt tokens served from cache
//   - OutputTokens: completion tokens (includes reasoning)
type codexRawTokenUsage struct {
	InputTokens       int `json:"input_tokens"`
	CachedInputTokens int `json:"cached_input_tokens"`
	OutputTokens      int `json:"output_tokens"`
}
