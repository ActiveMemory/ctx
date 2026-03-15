//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/method"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	session2 "github.com/ActiveMemory/ctx/internal/mcp/session"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Server is an MCP server that exposes ctx context over JSON-RPC 2.0.
//
// It reads JSON-RPC requests from stdin and writes responses to stdout,
// following the Model Context Protocol specification.
//
// Thread-safety: outMu serialises all writes to out (main loop + poller
// goroutine). The main loop itself is single-threaded, so request
// dispatch and session mutations need no additional locking.
type Server struct {
	contextDir  string
	version     string
	tokenBudget int
	out         io.Writer
	outMu       sync.Mutex // guards all writes to out
	in          io.Reader
	session     *session2.State
	poller      *ResourcePoller
}

// NewServer creates a new MCP server for the given context directory.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - version: binary version string for the server info response
//
// Returns:
//   - *Server: a configured MCP server ready to serve
func NewServer(contextDir, version string) *Server {
	srv := &Server{
		contextDir:  contextDir,
		version:     version,
		tokenBudget: rc.TokenBudget(),
		out:         os.Stdout,
		in:          os.Stdin,
		session:     session2.NewState(contextDir),
	}
	srv.poller = NewResourcePoller(contextDir, srv.emitNotification)
	return srv
}

// Serve starts the MCP server, reading from stdin and writing to stdout.
//
// It blocks until stdin is closed or an unrecoverable error occurs.
// Each line from stdin is expected to be a JSON-RPC 2.0 request.
//
// Returns:
//   - error: non-nil if an I/O error prevents continued operation
func (s *Server) Serve() error {
	defer s.poller.Stop()

	scanner := bufio.NewScanner(s.in)

	scanner.Buffer(make([]byte, 0, cfg.ScanMaxSize), cfg.ScanMaxSize)

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		resp := s.handleMessage(line)
		if resp == nil {
			// Notification: no response required.
			continue
		}

		out, err := json.Marshal(resp)
		if err != nil {
			// Marshal failure is an internal error; try to report it.
			s.writeError(nil, proto.ErrCodeInternal, assets.TextDesc(
				assets.TextDescKeyMCPFailedMarshal),
			)
			continue
		}
		s.outMu.Lock()
		_, writeErr := s.out.Write(append(out, token.NewlineLF[0]))
		s.outMu.Unlock()
		if writeErr != nil {
			return writeErr
		}
	}

	return scanner.Err()
}

// emitNotification writes a JSON-RPC notification to stdout.
// Safe to call from any goroutine (e.g., the resource poller).
func (s *Server) emitNotification(n proto.Notification) {
	out, err := json.Marshal(n)
	if err != nil {
		return
	}
	s.outMu.Lock()
	_, _ = s.out.Write(append(out, token.NewlineLF[0]))
	s.outMu.Unlock()
}

// handleMessage dispatches a raw JSON-RPC message to the appropriate
// handler.
//
// Parameters:
//   - data: raw JSON bytes from stdin
//
// Returns:
//   - *Response: JSON-RPC response, or nil for notifications
func (s *Server) handleMessage(data []byte) *proto.Response {
	var req proto.Request
	if err := json.Unmarshal(data, &req); err != nil {
		return &proto.Response{
			JSONRPC: server.JSONRPCVersion,
			Error: &proto.RPCError{
				Code:    proto.ErrCodeParse,
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
func (s *Server) dispatch(req proto.Request) *proto.Response {
	switch req.Method {
	case method.Initialize:
		return s.handleInitialize(req)
	case method.Ping:
		return s.ok(req.ID, struct{}{})
	case method.ResourcesList:
		return s.handleResourcesList(req)
	case method.ResourcesRead:
		return s.handleResourcesRead(req)
	case method.ResourcesSubscribe:
		return s.handleResourcesSubscribe(req)
	case method.ResourcesUnsubscribe:
		return s.handleResourcesUnsubscribe(req)
	case method.ToolsList:
		return s.handleToolsList(req)
	case method.ToolsCall:
		return s.handleToolsCall(req)
	case method.PromptsList:
		return s.handlePromptsList(req)
	case method.PromptsGet:
		return s.handlePromptsGet(req)
	default:
		return s.error(req.ID, proto.ErrCodeNotFound,
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
func (s *Server) handleNotification(req proto.Request) {
}

// handleInitialize responds to the MCP initialize handshake.
//
// Parameters:
//   - req: parsed JSON-RPC request
//
// Returns:
//   - *Response: server capabilities and protocol version
func (s *Server) handleInitialize(req proto.Request) *proto.Response {
	result := proto.InitializeResult{
		ProtocolVersion: proto.ProtocolVersion,
		Capabilities: proto.ServerCaps{
			Resources: &proto.ResourcesCap{Subscribe: true},
			Tools:     &proto.ToolsCap{},
			Prompts:   &proto.PromptsCap{},
		},
		ServerInfo: proto.AppInfo{
			Name:    server.Name,
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
func (s *Server) ok(id json.RawMessage, result interface{}) *proto.Response {
	return &proto.Response{
		JSONRPC: server.JSONRPCVersion,
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
func (s *Server) error(id json.RawMessage, code int, msg string) *proto.Response {
	return &proto.Response{
		JSONRPC: server.JSONRPCVersion,
		ID:      id,
		Error:   &proto.RPCError{Code: code, Message: msg},
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
