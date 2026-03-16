//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package out

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// OkResponse builds a successful JSON-RPC response.
//
// Parameters:
//   - id: request ID to echo back
//   - result: response payload
//
// Returns:
//   - *proto.Response: success response
func OkResponse(id json.RawMessage, result interface{}) *proto.Response {
	return &proto.Response{
		JSONRPC: server.JSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

// ErrResponse builds a JSON-RPC error response.
//
// Parameters:
//   - id: request ID to echo back
//   - code: JSON-RPC error code
//   - msg: human-readable error message
//
// Returns:
//   - *proto.Response: error response
func ErrResponse(id json.RawMessage, code int, msg string) *proto.Response {
	return &proto.Response{
		JSONRPC: server.JSONRPCVersion,
		ID:      id,
		Error:   &proto.RPCError{Code: code, Message: msg},
	}
}

// ToolOK builds a successful tool result with text content.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - text: success text to include in the result
//
// Returns:
//   - *proto.Response: tool result with text content
func ToolOK(id json.RawMessage, text string) *proto.Response {
	return OkResponse(
		id,
		proto.CallToolResult{
			Content: []proto.ToolContent{
				{Type: mime.ContentTypeText, Text: text},
			},
		})
}

// ToolError builds a tool error result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - msg: error message text
//
// Returns:
//   - *proto.Response: tool result with IsError set
func ToolError(id json.RawMessage, msg string) *proto.Response {
	return OkResponse(id, proto.CallToolResult{
		Content: []proto.ToolContent{{Type: mime.ContentTypeText, Text: msg}},
		IsError: true,
	})
}

// ToolResult wraps a handler (string, error) return into a
// proto.Response.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - text: success text from the handler
//   - err: handler error, nil on success
//
// Returns:
//   - *proto.Response: tool OK or tool error response
func ToolResult(
	id json.RawMessage, text string, err error,
) *proto.Response {
	if err != nil {
		return ToolError(id, err.Error())
	}
	return ToolOK(id, text)
}

// Call invokes a no-arg handler and wraps the result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - fn: handler function returning (string, error)
//
// Returns:
//   - *proto.Response: wrapped handler result
func Call(
	id json.RawMessage, fn func() (string, error),
) *proto.Response {
	text, err := fn()
	return ToolResult(id, text, err)
}
