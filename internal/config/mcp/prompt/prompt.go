//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

// MCP prompt name constants.
const (
	// SessionStart is the MCP prompt name for session initialization.
	SessionStart = "ctx-session-start"
	// AddDecision is the MCP prompt name for recording decisions.
	AddDecision = "ctx-decision-add"
	// AddLearning is the MCP prompt name for recording learnings.
	AddLearning = "ctx-learning-add"
	// Reflect is the MCP prompt name for session reflection.
	Reflect = "ctx-reflect"
	// Checkpoint is the MCP prompt name for session checkpoint.
	Checkpoint = "ctx-checkpoint"

	// RoleUser is the MCP message role for user-originated prompts.
	RoleUser = "user"
)
