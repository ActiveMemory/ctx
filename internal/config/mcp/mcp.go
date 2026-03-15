//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

// MCP constants.
const (
	// MCPResourceURIPrefix is the URI scheme prefix for MCP context resources.
	MCPResourceURIPrefix = "ctx://context/"
	// MimeMarkdown is the MIME type for Markdown content.
	MimeMarkdown = "text/markdown"
	// MCPScanMaxSize is the maximum scanner buffer size for MCP messages (1 MB).
	MCPScanMaxSize = 1 << 20
	// MCPMethodInitialize is the MCP initialize handshake method.
	MCPMethodInitialize = "initialize"
	// MCPMethodPing is the MCP ping method.
	MCPMethodPing = "ping"
	// MCPMethodResourcesList is the MCP method for listing resources.
	MCPMethodResourcesList = "resources/list"
	// MCPMethodResourcesRead is the MCP method for reading a resource.
	MCPMethodResourcesRead = "resources/read"
	// MCPMethodResourcesSubscribe is the MCP method for subscribing to resource changes.
	MCPMethodResourcesSubscribe = "resources/subscribe"
	// MCPMethodResourcesUnsubscribe is the MCP method for unsubscribing from resource changes.
	MCPMethodResourcesUnsubscribe = "resources/unsubscribe"
	// MCPMethodToolsList is the MCP method for listing tools.
	MCPMethodToolsList = "tools/list"
	// MCPMethodToolsCall is the MCP method for calling a tool.
	MCPMethodToolsCall = "tools/call"
	// MCPMethodPromptsList is the MCP method for listing prompts.
	MCPMethodPromptsList = "prompts/list"
	// MCPMethodPromptsGet is the MCP method for getting a prompt.
	MCPMethodPromptsGet = "prompts/get"
	// MCPJSONRPCVersion is the JSON-RPC protocol version string.
	MCPJSONRPCVersion = "2.0"
	// MCPServerName is the server name reported during initialization.
	MCPServerName = "ctx"
	// MCPContentTypeText is the content type for text tool output.
	MCPContentTypeText = "text"
	// MCPSchemaObject is the JSON Schema type for objects.
	MCPSchemaObject = "object"
	// MCPSchemaString is the JSON Schema type for strings.
	MCPSchemaString = "string"
	// MCPSchemaNumber is the JSON Schema type for numbers.
	MCPSchemaNumber = "number"
	// MCPSchemaBoolean is the JSON Schema type for booleans.
	MCPSchemaBoolean = "boolean"
	// MCPToolStatus is the MCP tool name for context status.
	MCPToolStatus = "ctx_status"
	// MCPToolAdd is the MCP tool name for adding entries.
	MCPToolAdd = "ctx_add"
	// MCPToolComplete is the MCP tool name for completing tasks.
	MCPToolComplete = "ctx_complete"
	// MCPToolDrift is the MCP tool name for drift detection.
	MCPToolDrift = "ctx_drift"
	// MCPToolRecall is the MCP tool name for querying session history.
	MCPToolRecall = "ctx_recall"
	// MCPToolWatchUpdate is the MCP tool name for structured context updates.
	MCPToolWatchUpdate = "ctx_watch_update"
	// MCPToolCompact is the MCP tool name for compacting tasks.
	MCPToolCompact = "ctx_compact"
	// MCPToolNext is the MCP tool name for suggesting the next task.
	MCPToolNext = "ctx_next"
	// MCPToolCheckTaskCompletion is the MCP tool name for task completion nudge.
	MCPToolCheckTaskCompletion = "ctx_check_task_completion"
	// MCPToolSessionEvent is the MCP tool name for session lifecycle events.
	MCPToolSessionEvent = "ctx_session_event"
	// MCPToolRemind is the MCP tool name for listing reminders.
	MCPToolRemind = "ctx_remind"

	// MCPPromptSessionStart is the MCP prompt name for session initialization.
	MCPPromptSessionStart = "ctx-session-start"
	// MCPPromptAddDecision is the MCP prompt name for recording decisions.
	MCPPromptAddDecision = "ctx-add-decision"
	// MCPPromptAddLearning is the MCP prompt name for recording learnings.
	MCPPromptAddLearning = "ctx-add-learning"
	// MCPPromptReflect is the MCP prompt name for session reflection.
	MCPPromptReflect = "ctx-reflect"
	// MCPPromptCheckpoint is the MCP prompt name for session checkpoint.
	MCPPromptCheckpoint = "ctx-checkpoint"

	// MCPNotifyResourcesUpdated is the MCP notification for resource changes.
	MCPNotifyResourcesUpdated = "notifications/resources/updated"
)
