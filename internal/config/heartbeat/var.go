//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package heartbeat

// Template variable keys for heartbeat payloads.
const (
	// VarPromptCount is the heartbeat field for prompt count.
	VarPromptCount = "prompt_count"
	// VarSessionID is the heartbeat field for session identifier.
	VarSessionID = "session_id"
	// VarContextModified is the heartbeat field for context modification flag.
	VarContextModified = "context_modified"
	// VarTokens is the heartbeat field for token count.
	VarTokens = "tokens"
	// VarContextWindow is the heartbeat field for context window size.
	VarContextWindow = "context_window"
	// VarUsagePct is the heartbeat field for usage percentage.
	VarUsagePct = "usage_pct"
)
