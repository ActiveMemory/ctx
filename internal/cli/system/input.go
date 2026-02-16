//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"encoding/json"
	"io"
)

// HookInput represents the JSON payload that Claude Code sends to hook
// commands via stdin.
type HookInput struct {
	SessionID string    `json:"session_id"`
	ToolInput ToolInput `json:"tool_input"`
}

// ToolInput contains the tool-specific fields from a Claude Code hook
// invocation. For Bash hooks, Command holds the shell command.
type ToolInput struct {
	Command string `json:"command"`
}

// readInput reads and parses the JSON hook input from r.
// Returns a zero-value HookInput on any error (graceful degradation).
func readInput(r io.Reader) HookInput {
	var input HookInput
	data, err := io.ReadAll(r)
	if err != nil {
		return input
	}
	_ = json.Unmarshal(data, &input)
	return input
}
