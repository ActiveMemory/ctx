//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package proto

import "encoding/json"

// JSON-RPC 2.0 message types for the Model Context Protocol.
//
// See: https://spec.modelcontextprotocol.io/

// Request represents a JSON-RPC 2.0 request message.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response message.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// Notification represents a JSON-RPC 2.0 notification (no id).
type Notification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error object.
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes.
const (
	ErrCodeParse      = -32700
	errCodeInvalidReq = -32600
	ErrCodeNotFound   = -32601
	ErrCodeInvalidArg = -32602
	ErrCodeInternal   = -32603
)

// ProtocolVersion is the MCP protocol version.
const ProtocolVersion = "2024-11-05"

// --- Initialization types ---

// InitializeParams contains client information sent during initialization.
type InitializeParams struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    ClientCaps `json:"capabilities"`
	ClientInfo      AppInfo    `json:"clientInfo"`
}

// ClientCaps describes client capabilities.
type ClientCaps struct {
	Roots    *struct{} `json:"roots,omitempty"`
	Sampling *struct{} `json:"sampling,omitempty"`
}

// AppInfo identifies a client or server application.
type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is the server's response to initialize.
type InitializeResult struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    ServerCaps `json:"capabilities"`
	ServerInfo      AppInfo    `json:"serverInfo"`
}

// ServerCaps describes server capabilities.
type ServerCaps struct {
	Resources *ResourcesCap `json:"resources,omitempty"`
	Tools     *ToolsCap     `json:"tools,omitempty"`
	Prompts   *PromptsCap   `json:"prompts,omitempty"`
}

// ResourcesCap indicates the server supports resources.
type ResourcesCap struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

// ToolsCap indicates the server supports tools.
type ToolsCap struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// --- Resource types ---

// Resource describes a single MCP resource.
type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ResourceListResult is returned by resources/list.
type ResourceListResult struct {
	Resources []Resource `json:"resources"`
}

// ReadResourceParams is sent with resources/read.
type ReadResourceParams struct {
	URI string `json:"uri"`
}

// ResourceContent represents the content of a resource.
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
}

// ReadResourceResult is returned by resources/read.
type ReadResourceResult struct {
	Contents []ResourceContent `json:"contents"`
}

// SubscribeParams is sent with resources/subscribe.
type SubscribeParams struct {
	URI string `json:"uri"`
}

// UnsubscribeParams is sent with resources/unsubscribe.
type UnsubscribeParams struct {
	URI string `json:"uri"`
}

// ResourceUpdatedParams is sent with notifications/resources/updated.
type ResourceUpdatedParams struct {
	URI string `json:"uri"`
}

// --- Tool types ---

// ToolAnnotations provides hints about a tool's behavior.
type ToolAnnotations struct {
	ReadOnlyHint    bool `json:"readOnlyHint,omitempty"`
	DestructiveHint bool `json:"destructiveHint,omitempty"`
	IdempotentHint  bool `json:"idempotentHint,omitempty"`
}

// Tool describes a single MCP tool.
type Tool struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	InputSchema InputSchema      `json:"inputSchema"`
	Annotations *ToolAnnotations `json:"annotations,omitempty"`
}

// InputSchema describes the JSON Schema for tool inputs.
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property describes a single property in a JSON Schema.
type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolListResult is returned by tools/list.
type ToolListResult struct {
	Tools []Tool `json:"tools"`
}

// CallToolParams is sent with tools/call.
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolContent represents a piece of tool output.
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// CallToolResult is returned by tools/call.
type CallToolResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// --- Prompt types ---

// PromptsCap indicates the server supports prompts.
type PromptsCap struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// Prompt describes a single MCP prompt template.
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument describes a single argument for a prompt.
type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// PromptListResult is returned by prompts/list.
type PromptListResult struct {
	Prompts []Prompt `json:"prompts"`
}

// GetPromptParams is sent with prompts/get.
type GetPromptParams struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

// PromptMessage represents a message in a prompt response.
type PromptMessage struct {
	Role    string      `json:"role"`
	Content ToolContent `json:"content"`
}

// GetPromptResult is returned by prompts/get.
type GetPromptResult struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}
