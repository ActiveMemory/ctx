//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

// response is the JSON output format for Claude Code hooks.
// Using structured JSON ensures the agent processes the output as a directive
// rather than treating it as ignorable plain text.
type response struct {
	Output *responseOutput `json:"hookSpecificOutput,omitempty"`
}

// responseOutput carries event-specific fields inside a response.
type responseOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext,omitempty"`
}

// usageData represents the minimal usage fields from a Claude Code JSONL
// assistant message. Only the fields needed for token counting are included.
type usageData struct {
	InputTokens              int `json:"input_tokens"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens"`
}

// jsonlMessage represents the minimal structure of a Claude Code JSONL line
// needed to extract usage and model data from assistant messages.
type jsonlMessage struct {
	Type    string `json:"type"`
	Message struct {
		Role  string    `json:"role"`
		Model string    `json:"model"`
		Usage usageData `json:"usage"`
	} `json:"message"`
}
