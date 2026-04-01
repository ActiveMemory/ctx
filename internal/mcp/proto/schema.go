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
//
// Fields:
//   - JSONRPC: Protocol version, always "2.0"
//   - ID: Request identifier for correlating responses
//   - Method: RPC method name
//   - Params: Method-specific parameters
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response message.
//
// Fields:
//   - JSONRPC: Protocol version, always "2.0"
//   - ID: Request identifier this response correlates to
//   - Result: Success payload (mutually exclusive with Error)
//   - Error: Error payload (mutually exclusive with Result)
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

// Notification represents a JSON-RPC 2.0 notification (no id).
//
// Fields:
//   - JSONRPC: Protocol version, always "2.0"
//   - Method: Notification method name
//   - Params: Method-specific parameters
type Notification struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RPCError represents a JSON-RPC 2.0 error object.
//
// Fields:
//   - Code: Numeric error code (see ErrCode* constants)
//   - Message: Human-readable error description
//   - Data: Optional structured error data
type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes.
const (
	ErrCodeParse      = -32700
	ErrCodeInvalidReq = -32600
	ErrCodeNotFound   = -32601
	ErrCodeInvalidArg = -32602
	ErrCodeInternal   = -32603
)

// ProtocolVersion is the MCP protocol version.
const ProtocolVersion = "2024-11-05"

// --- Initialization types ---

// InitializeParams contains client information sent during initialization.
//
// Fields:
//   - ProtocolVersion: MCP version the client supports
//   - Capabilities: Client capability flags
//   - ClientInfo: Client name and version
type InitializeParams struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    ClientCaps `json:"capabilities"`
	ClientInfo      AppInfo    `json:"clientInfo"`
}

// ClientCaps describes client capabilities.
//
// Fields:
//   - Roots: Non-nil if client supports roots
//   - Sampling: Non-nil if client supports sampling
type ClientCaps struct {
	Roots    *struct{} `json:"roots,omitempty"`
	Sampling *struct{} `json:"sampling,omitempty"`
}

// AppInfo identifies a client or server application.
//
// Fields:
//   - Name: Application name
//   - Version: Application version string
type AppInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is the server's response to initialize.
//
// Fields:
//   - ProtocolVersion: MCP version the server chose
//   - Capabilities: Server capability flags
//   - ServerInfo: Server name and version
type InitializeResult struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    ServerCaps `json:"capabilities"`
	ServerInfo      AppInfo    `json:"serverInfo"`
}

// ServerCaps describes server capabilities.
//
// Fields:
//   - Resources: Non-nil if server supports resources
//   - Tools: Non-nil if server supports tools
//   - Prompts: Non-nil if server supports prompts
type ServerCaps struct {
	Resources *ResourcesCap `json:"resources,omitempty"`
	Tools     *ToolsCap     `json:"tools,omitempty"`
	Prompts   *PromptsCap   `json:"prompts,omitempty"`
}

// ResourcesCap indicates the server supports resources.
//
// Fields:
//   - Subscribe: Whether clients can subscribe to changes
//   - ListChanged: Whether the resource list can change
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
//
// Fields:
//   - URI: Unique resource identifier
//   - Name: Human-readable name
//   - Description: What the resource contains
//   - MimeType: Content type hint
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
//
// Fields:
//   - URI: Resource identifier
//   - MimeType: Content type
//   - Text: Text content
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
//
// Fields:
//   - ReadOnlyHint: Tool does not modify state
//   - DestructiveHint: Tool may cause irreversible changes
//   - IdempotentHint: Repeated calls produce the same result
type ToolAnnotations struct {
	ReadOnlyHint    bool `json:"readOnlyHint,omitempty"`
	DestructiveHint bool `json:"destructiveHint,omitempty"`
	IdempotentHint  bool `json:"idempotentHint,omitempty"`
}

// Tool describes a single MCP tool.
//
// Fields:
//   - Name: Tool identifier
//   - Description: What the tool does
//   - InputSchema: JSON Schema for tool arguments
//   - Annotations: Optional behavioral hints
type Tool struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	InputSchema InputSchema      `json:"inputSchema"`
	Annotations *ToolAnnotations `json:"annotations,omitempty"`
}

// InputSchema describes the JSON Schema for tool inputs.
//
// Fields:
//   - Type: Schema type, always "object"
//   - Properties: Named property definitions
//   - Required: List of required property names
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property describes a single property in a JSON Schema.
//
// Fields:
//   - Type: JSON type (string, integer, boolean, etc.)
//   - Description: Human-readable property description
//   - Enum: Allowed values (optional)
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
//
// Fields:
//   - Name: Tool to invoke
//   - Arguments: Tool-specific arguments
type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolContent represents a piece of tool output.
//
// Fields:
//   - Type: Content type (e.g. "text")
//   - Text: Text content
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// CallToolResult is returned by tools/call.
//
// Fields:
//   - Content: Output content pieces
//   - IsError: Whether the tool invocation failed
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
//
// Fields:
//   - Name: Prompt identifier
//   - Description: What the prompt does
//   - Arguments: Accepted arguments
type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

// PromptArgument describes a single argument for a prompt.
//
// Fields:
//   - Name: Argument name
//   - Description: What the argument is for
//   - Required: Whether the argument must be provided
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
//
// Fields:
//   - Name: Prompt to retrieve
//   - Arguments: Values for prompt arguments
type GetPromptParams struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments,omitempty"`
}

// PromptMessage represents a message in a prompt response.
//
// Fields:
//   - Role: Message role (user, assistant)
//   - Content: Message content
type PromptMessage struct {
	Role    string      `json:"role"`
	Content ToolContent `json:"content"`
}

// GetPromptResult is returned by prompts/get.
//
// Fields:
//   - Description: Prompt description
//   - Messages: Rendered prompt messages
type GetPromptResult struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}
