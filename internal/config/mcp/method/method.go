//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

const (
	// MethodInitialize is the MCP initialize handshake method.
	MethodInitialize = "initialize"
	// MethodPing is the MCP ping method.
	MethodPing = "ping"
	// MethodResourcesList is the MCP method for listing resources.
	MethodResourcesList = "resources/list"
	// MethodResourcesRead is the MCP method for reading a resource.
	MethodResourcesRead = "resources/read"
	// MethodResourcesSubscribe is the MCP method for subscribing to resource changes.
	MethodResourcesSubscribe = "resources/subscribe"
	// MethodResourcesUnsubscribe is the MCP method for unsubscribing from resource changes.
	MethodResourcesUnsubscribe = "resources/unsubscribe"
	// MethodToolsList is the MCP method for listing tools.
	MethodToolsList = "tools/list"
	// MethodToolsCall is the MCP method for calling a tool.
	MethodToolsCall = "tools/call"
	// MethodPromptsList is the MCP method for listing prompts.
	MethodPromptsList = "prompts/list"
	// MethodPromptsGet is the MCP method for getting a prompt.
	MethodPromptsGet = "prompts/get"
)
