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

// Claude API message roles.
const (
	// RoleUser is a user message.
	RoleUser = "user"
	// RoleAssistant is an assistant message.
	RoleAssistant = "assistant"
)
