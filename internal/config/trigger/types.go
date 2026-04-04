//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trigger

// TriggerType identifies the lifecycle trigger event type.
type TriggerType = string

// Lifecycle trigger event type constants.
const (
	// PreToolUse fires before an AI tool invocation.
	PreToolUse TriggerType = "pre-tool-use"
	// PostToolUse fires after an AI tool invocation.
	PostToolUse TriggerType = "post-tool-use"
	// SessionStart fires when an AI session begins.
	SessionStart TriggerType = "session-start"
	// SessionEnd fires when an AI session ends.
	SessionEnd TriggerType = "session-end"
	// FileSave fires when a file is saved.
	FileSave TriggerType = "file-save"
	// ContextAdd fires when context is added.
	ContextAdd TriggerType = "context-add"
)
