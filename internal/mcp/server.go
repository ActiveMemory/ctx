//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ActiveMemory/ctx/internal/config"
)

// Server is an MCP server that exposes ctx context over JSON-RPC 2.0.
//
// It reads JSON-RPC requests from stdin and writes responses to stdout,
// following the Model Context Protocol specification.
type Server struct {
	contextDir string
	version    string
	out        io.Writer
	in         io.Reader
}

// NewServer creates a new MCP server for the given context directory.
//
// Parameters:
//   - contextDir: Path to the .context/ directory
//
// Returns:
//   - *Server: A configured MCP server ready to serve
func NewServer(contextDir string) *Server {
	return &Server{
		contextDir: contextDir,
		version:    config.BinaryVersion,
		out:        os.Stdout,
		in:         os.Stdin,
	}
}

// Serve starts the MCP server, reading from stdin and writing to stdout.
//
// It blocks until stdin is closed or an unrecoverable error occurs.
// Each line from stdin is expected to be a JSON-RPC 2.0 request.
//
// Returns:
//   - error: Non-nil if an I/O error prevents continued operation
func (s *Server) Serve() error {
	scanner := bufio.NewScanner(s.in)

	// Increase scanner buffer for large messages (1MB).
	const maxScanSize = 1 << 20
	scanner.Buffer(make([]byte, 0, maxScanSize), maxScanSize)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		resp := s.handleMessage(line)
		if resp == nil {
			// Notification — no response required.
			continue
		}

		out, err := json.Marshal(resp)
		if err != nil {
			// Marshal failure is an internal error; try to report it.
			s.writeError(nil, errCodeInternal, "failed to marshal response")
			continue
		}
		_, _ = fmt.Fprintf(s.out, "%s\n", out) //nolint:gosec // G705: stdout, not HTTP response
	}

	return scanner.Err()
}

// handleMessage dispatches a raw JSON-RPC message to the appropriate handler.
func (s *Server) handleMessage(data []byte) *Response {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return &Response{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: errCodeParse, Message: "parse error"},
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
func (s *Server) dispatch(req Request) *Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "ping":
		return s.ok(req.ID, struct{}{})
	case "resources/list":
		return s.handleResourcesList(req)
	case "resources/read":
		return s.handleResourcesRead(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	default:
		return s.error(req.ID, errCodeNotFound,
			fmt.Sprintf("method not found: %s", req.Method))
	}
}

// handleNotification processes notifications (no response needed).
func (s *Server) handleNotification(req Request) {
	// MCP notifications we handle:
	// - notifications/initialized: client confirms init complete
	// - notifications/cancelled: client cancels a request
	// All are no-ops for our stateless server.
}

// handleInitialize responds to the MCP initialize handshake.
func (s *Server) handleInitialize(req Request) *Response {
	result := InitializeResult{
		ProtocolVersion: protocolVersion,
		Capabilities: ServerCaps{
			Resources: &ResourcesCap{},
			Tools:     &ToolsCap{},
		},
		ServerInfo: AppInfo{
			Name:    "ctx",
			Version: s.version,
		},
	}
	return s.ok(req.ID, result)
}

// ok builds a successful JSON-RPC response.
func (s *Server) ok(id json.RawMessage, result interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// error builds a JSON-RPC error response.
func (s *Server) error(id json.RawMessage, code int, msg string) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &RPCError{Code: code, Message: msg},
	}
}

// writeError writes an error response directly to stdout. Used when the
// normal response flow cannot be used (e.g., marshal failure).
func (s *Server) writeError(id json.RawMessage, code int, msg string) {
	resp := s.error(id, code, msg)
	if out, err := json.Marshal(resp); err == nil {
		_, _ = fmt.Fprintf(s.out, "%s\n", out) //nolint:gosec // G705: stdout, not HTTP response
	}
}
