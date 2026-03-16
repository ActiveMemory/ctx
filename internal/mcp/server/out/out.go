//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package out

import (
	"encoding/json"

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
