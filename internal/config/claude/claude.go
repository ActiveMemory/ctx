//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

// Claude model identification.
const (
	// ModelPrefix is the prefix for all Claude model IDs in JSONL.
	ModelPrefix = "claude-"
	// ModelSuffix1M is the suffix indicating a 1M context window model.
	ModelSuffix1M = "[1m]"
)

// Claude API content block types.
const (
	// BlockText is a text content block.
	BlockText = "text"
	// BlockThinking is an extended thinking content block.
	BlockThinking = "thinking"
	// BlockToolUse is a tool invocation block.
	BlockToolUse = "tool_use"
	// BlockToolResult is a tool execution result block.
	BlockToolResult = "tool_result"
)

// JSONL field names for quick-scan filtering.
const (
	// FieldUsage is the JSON key for the usage block.
	FieldUsage = "usage"
	// FieldInputTokens is the JSON key for input token count.
	FieldInputTokens = "input_tokens"
)

// Model family substrings for context window detection.
const (
	// ModelOpus is the substring identifying Opus models (always 1M).
	ModelOpus = "opus"
)

// Context window sizes.
const (
	// ContextWindow1M is the context window for 1M-capable models.
	ContextWindow1M = 1_000_000
)

// JSONL file scanning limits.
const (
	// MaxTailBytes is the max bytes to read from the end of a JSONL file.
	MaxTailBytes = 32768
)

// Content block detection prefixes.
const (
	// ContentBlockArrayPrefix is the JSON prefix for Claude's content
	// block array format ([{type: "text", text: "..."}]).
	ContentBlockArrayPrefix = "[{"
)

// Claude API message roles.
const (
	// RoleUser is a user message.
	RoleUser = "user"
	// RoleAssistant is an assistant message.
	RoleAssistant = "assistant"
)
