//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tpl

// Template variable key constants — used as map keys in template.Execute
// data maps to avoid magic strings in hook and display code.
const (
	// VarAlertMessages is the template variable for resource alert messages.
	VarAlertMessages = "AlertMessages"
	// VarUnenrichedCount is the template variable for unenriched entry count.
	VarUnenrichedCount = "UnenrichedCount"
	// VarUnexportedCount is the template variable for unexported session count.
	VarUnexportedCount = "UnexportedCount"
	// VarBinaryVersion is the template variable for the binary version string.
	VarBinaryVersion = "BinaryVersion"
	// VarFileWarnings is the template variable for knowledge file warnings.
	VarFileWarnings = "FileWarnings"
	// VarKeyAgeDays is the template variable for API key age in days.
	VarKeyAgeDays = "KeyAgeDays"
	// VarLastRefreshDate is the template variable for the last map refresh date.
	VarLastRefreshDate = "LastRefreshDate"
	// VarModuleCount is the template variable for the number of changed modules.
	VarModuleCount = "ModuleCount"
	// VarPercentage is the template variable for context window percentage.
	VarPercentage = "Percentage"
	// VarPluginVersion is the template variable for the plugin version string.
	VarPluginVersion = "PluginVersion"
	// VarPromptCount is the template variable for the prompt counter.
	VarPromptCount = "PromptCount"
	// VarPromptsSinceNudge is the template variable for prompts since last nudge.
	VarPromptsSinceNudge = "PromptsSinceNudge"
	// VarReminderList is the template variable for formatted reminder list.
	VarReminderList = "ReminderList"
	// VarThreshold is the template variable for a token threshold value.
	VarThreshold = "Threshold"
	// VarTokenCount is the template variable for a token count value.
	VarTokenCount = "TokenCount"
	// VarWarnings is the template variable for backup warning messages.
	VarWarnings = "Warnings"
	// VarHeartbeatPromptCount is the heartbeat field for prompt count.
	VarHeartbeatPromptCount = "prompt_count"
	// VarHeartbeatSessionID is the heartbeat field for session identifier.
	VarHeartbeatSessionID = "session_id"
	// VarHeartbeatContextModified is the heartbeat field for context modification flag.
	VarHeartbeatContextModified = "context_modified"
	// VarHeartbeatTokens is the heartbeat field for token count.
	VarHeartbeatTokens = "tokens"
	// VarHeartbeatContextWindow is the heartbeat field for context window size.
	VarHeartbeatContextWindow = "context_window"
	// VarHeartbeatUsagePct is the heartbeat field for usage percentage.
	VarHeartbeatUsagePct = "usage_pct"
)
