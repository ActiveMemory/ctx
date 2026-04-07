//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

// JSONL top-level field names for message records.

// Required fields present on every user/assistant record.
var RequiredFields = []string{
	"uuid", "parentUuid", "sessionId", "timestamp",
	"type", "cwd", "version", "message",
	"isSidechain", "userType",
}

// Optional fields that may appear on message records.
var OptionalFields = []string{
	"gitBranch", "slug", "requestId",
	"thinkingMetadata", "todos", "permissionMode",
	"logicalParentUuid", "isMeta", "compactMetadata",
	"isVisibleInTranscriptOnly", "isCompactSummary",
	"agentId", "teamName", "agentName", "agentColor",
	"promptId", "entrypoint", "agentSetting",
	"sourceToolAssistantUUID", "toolUseResult",
	"sourceToolUseID", "origin", "planContent",
	"isApiErrorMessage", "error", "apiError",
}

// JSONL record type values.

// Message record types that carry conversation content.
const (
	// RecordUser is a user message record.
	RecordUser = "user"
	// RecordAssistant is an assistant message record.
	RecordAssistant = "assistant"
)

// Metadata record types (no field validation).
var MetadataRecordTypes = []string{
	"last-prompt", "custom-title", "ai-title",
	"attachment", "permission-mode", "agent-name",
	"agent-color", "agent-setting", "tag", "pr-link",
	"mode", "worktree-state", "content-replacement",
	"speculation-accept", "task-summary",
}

// Infrastructure record types (skipped by parser).
var InfraRecordTypes = []string{
	"progress", "file-history-snapshot",
	"attribution-snapshot", "system",
	"summary", "queue-operation",
}

// Content block type values.

// Parsed block types that the journal parser extracts.
var ParsedBlockTypes = []string{
	"text", "thinking", "tool_use", "tool_result",
}

// Known block types recognized but not parsed.
var KnownBlockTypes = []string{
	"server_tool_use", "mcp_tool_use",
	"mcp_tool_result", "code_execution_tool_result",
	"container_upload",
}
