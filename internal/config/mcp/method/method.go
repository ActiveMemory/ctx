//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package method

const (
	// Initialize is the MCP initialize handshake method.
	Initialize = "initialize"
	// Ping is the MCP ping method.
	Ping = "ping"
	// ResourcesList is the MCP method for listing resources.
	ResourcesList = "resources/list"
	// ResourcesRead is the MCP method for reading a resource.
	ResourcesRead = "resources/read"
	// ResourcesSubscribe is the MCP method for subscribing to resource changes.
	ResourcesSubscribe = "resources/subscribe"
	// ResourcesUnsubscribe is the MCP method for unsubscribing from resource
	// changes.
	ResourcesUnsubscribe = "resources/unsubscribe"
	// ToolsList is the MCP method for listing tools.
	ToolsList = "tools/list"
	// ToolsCall is the MCP method for calling a tool.
	ToolsCall = "tools/call"
	// PromptsList is the MCP method for listing prompts.
	PromptsList = "prompts/list"
	// PromptsGet is the MCP method for getting a prompt.
	PromptsGet = "prompts/get"
)
