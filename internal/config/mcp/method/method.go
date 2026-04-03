//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package method

// MCP JSON-RPC method constants.
const (
	// Initialize is the MCP initialize handshake method.
	Initialize = "initialize"
	// Ping is the MCP ping method.
	Ping = "ping"
	// ResourceList is the MCP method for listing resources.
	ResourceList = "resources/list"
	// ResourceRead is the MCP method for reading a resource.
	ResourceRead = "resources/read"
	// ResourceSubscribe is the MCP method for subscribing to resource changes.
	ResourceSubscribe = "resources/subscribe"
	// ResourceUnsubscribe is the MCP method for unsubscribing from resource
	// changes.
	ResourceUnsubscribe = "resources/unsubscribe"
	// ToolList is the MCP method for listing tools.
	ToolList = "tools/list"
	// ToolCall is the MCP method for calling a tool.
	ToolCall = "tools/call"
	// PromptList is the MCP method for listing prompts.
	PromptList = "prompts/list"
	// PromptGet is the MCP method for getting a prompt.
	PromptGet = "prompts/get"
)
