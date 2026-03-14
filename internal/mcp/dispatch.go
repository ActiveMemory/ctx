//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// handleMessage dispatches a raw JSON-RPC message to the appropriate
// handler.
//
// Parameters:
//   - data: raw JSON bytes from stdin
//
// Returns:
//   - *Response: JSON-RPC response, or nil for notifications
func (s *Server) handleMessage(data []byte) *Response {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return &Response{
			JSONRPC: mcp.MCPJSONRPCVersion,
			Error: &RPCError{
				Code:    errCodeParse,
				Message: assets.TextDesc(assets.TextDescKeyMCPParseError),
			},
		}
	}

	// Notifications have no ID and expect no response.
	if req.ID == nil {
		s.handleNotification(req)
		return nil
	}

	return s.dispatch(req)
}

// dispatch routes a request to the correct handler based on method name.
//
// Parameters:
//   - req: parsed JSON-RPC request
//
// Returns:
//   - *Response: result or error response
func (s *Server) dispatch(req Request) *Response {
	switch req.Method {
	case mcp.MCPMethodInitialize:
		return s.handleInitialize(req)
	case mcp.MCPMethodPing:
		return s.ok(req.ID, struct{}{})
	case mcp.MCPMethodResourcesList:
		return s.handleResourcesList(req)
	case mcp.MCPMethodResourcesRead:
		return s.handleResourcesRead(req)
	case mcp.MCPMethodResourcesSubscribe:
		return s.handleResourcesSubscribe(req)
	case mcp.MCPMethodResourcesUnsubscribe:
		return s.handleResourcesUnsubscribe(req)
	case mcp.MCPMethodToolsList:
		return s.handleToolsList(req)
	case mcp.MCPMethodToolsCall:
		return s.handleToolsCall(req)
	case mcp.MCPMethodPromptsList:
		return s.handlePromptsList(req)
	case mcp.MCPMethodPromptsGet:
		return s.handlePromptsGet(req)
	default:
		return s.error(req.ID, errCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPMethodNotFound), req.Method),
		)
	}
}

// handleNotification processes notifications (no response needed).
//
// MCP notifications handled:
//   - notifications/initialized: client confirms init complete
//   - notifications/cancelled: client cancels a request
//
// All are no-ops for our stateless server.
//
// Parameters:
//   - req: parsed JSON-RPC notification
func (s *Server) handleNotification(req Request) {
}

// handleInitialize responds to the MCP initialize handshake.
//
// Parameters:
//   - req: parsed JSON-RPC request
//
// Returns:
//   - *Response: server capabilities and protocol version
func (s *Server) handleInitialize(req Request) *Response {
	result := InitializeResult{
		ProtocolVersion: protocolVersion,
		Capabilities: ServerCaps{
			Resources: &ResourcesCap{Subscribe: true},
			Tools:     &ToolsCap{},
			Prompts:   &PromptsCap{},
		},
		ServerInfo: AppInfo{
			Name:    mcp.MCPServerName,
			Version: s.version,
		},
	}
	return s.ok(req.ID, result)
}

// ok builds a successful JSON-RPC response.
//
// Parameters:
//   - id: request ID to echo back
//   - result: response payload
//
// Returns:
//   - *Response: success response
func (s *Server) ok(id json.RawMessage, result interface{}) *Response {
	return &Response{
		JSONRPC: mcp.MCPJSONRPCVersion,
		ID:      id,
		Result:  result,
	}
}

// error builds a JSON-RPC error response.
//
// Parameters:
//   - id: request ID to echo back
//   - code: JSON-RPC error code
//   - msg: human-readable error message
//
// Returns:
//   - *Response: error response
func (s *Server) error(id json.RawMessage, code int, msg string) *Response {
	return &Response{
		JSONRPC: mcp.MCPJSONRPCVersion,
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg},
	}
}

// writeError writes an error response directly to stdout.
//
// Used when the normal response flow cannot be used (e.g., marshal
// failure). This is a last-resort fallback; write failures are
// silently ignored.
//
// Parameters:
//   - id: request ID to echo back (may be nil)
//   - code: JSON-RPC error code
//   - msg: human-readable error message
func (s *Server) writeError(id json.RawMessage, code int, msg string) {
	resp := s.error(id, code, msg)
	if out, marshalErr := json.Marshal(resp); marshalErr == nil {
		s.outMu.Lock()
		_, _ = s.out.Write(append(out, token.NewlineLF[0]))
		s.outMu.Unlock()
	}
}
