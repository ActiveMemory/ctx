//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// TriggerType identifies the lifecycle trigger event type.
type TriggerType string

const (
	// TriggerPreToolUse fires before an AI tool invocation.
	TriggerPreToolUse TriggerType = "pre-tool-use"
	// TriggerPostToolUse fires after an AI tool invocation.
	TriggerPostToolUse TriggerType = "post-tool-use"
	// TriggerSessionStart fires when an AI session begins.
	TriggerSessionStart TriggerType = "session-start"
	// TriggerSessionEnd fires when an AI session ends.
	TriggerSessionEnd TriggerType = "session-end"
	// TriggerFileSave fires when a file is saved.
	TriggerFileSave TriggerType = "file-save"
	// TriggerContextAdd fires when context is added.
	TriggerContextAdd TriggerType = "context-add"
)

// TriggerInput is the JSON object sent to trigger scripts via stdin.
//
// Fields:
//   - TriggerType: Lifecycle event category
//   - Tool: Name of the AI tool being used
//   - Parameters: Tool-specific parameters
//   - Session: Session metadata (id and model)
//   - Timestamp: ISO 8601 timestamp
//   - CtxVersion: Version of ctx
type TriggerInput struct {
	TriggerType string         `json:"hookType"`
	Tool        string         `json:"tool"`
	Parameters  map[string]any `json:"parameters"`
	Session     TriggerSession `json:"session"`
	Timestamp   string         `json:"timestamp"`
	CtxVersion  string         `json:"ctxVersion"`
}

// TriggerSession contains session metadata sent to trigger scripts.
type TriggerSession struct {
	ID    string `json:"id"`
	Model string `json:"model"`
}
