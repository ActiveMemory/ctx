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
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	res "github.com/ActiveMemory/ctx/internal/mcp/server/resource"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Server is an MCP server that exposes ctx context over JSON-RPC 2.0.
//
// It reads JSON-RPC requests from stdin and writes responses to stdout,
// following the Model Context Protocol specification.
//
// Thread-safety: outMu serialises all writes to out (main loop and poller
// goroutine). The main loop itself is single-threaded, so request
// dispatch and session mutations need no additional locking.
type Server struct {
	handler      *handler.Handler
	version      string
	out          io.Writer
	outMu        sync.Mutex // guards all writes to out
	in           io.Reader
	poller       *ResourcePoller
	resourceList proto.ResourceListResult // pre-built, immutable after init
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
	res.Init()
	srv := &Server{
		handler:      handler.New(contextDir, rc.TokenBudget()),
		version:      version,
		out:          os.Stdout,
		in:           os.Stdin,
		resourceList: res.ToList(),
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

		if writeErr := s.writeJSON(resp); writeErr != nil {
			// Marshal failure: try to report it as an error response.
			fallback := out.ErrResponse(nil, proto.ErrCodeInternal,
				assets.TextDesc(assets.TextDescKeyMCPFailedMarshal))
			if fbErr := s.writeJSON(fallback); fbErr != nil {
				return fbErr
			}
			continue
		}
	}

	return scanner.Err()
}

// emitNotification writes a JSON-RPC notification to stdout.
//
// Safe to call from any goroutine (e.g., the resource poller).
// Write failures are silently ignored — notifications are best-effort.
//
// Parameters:
//   - n: notification to marshal and write
func (s *Server) emitNotification(n proto.Notification) {
	_ = s.writeJSON(n)
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

// dispatch routes a request to the correct handler based on the method name.
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
		return out.OkResponse(req.ID, struct{}{})
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
		return out.ErrResponse(req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPMethodNotFound), req.Method),
		)
	}
}

// handleNotification processes notifications (no response needed).
//
// MCP notifications handled:
//   - notifications/initialized: the client confirms init complete
//   - notifications/canceled: the client cancels a request
//
// All are no-ops for our stateless server.
//
// Parameters:
//   - _: parsed JSON-RPC notification
func (s *Server) handleNotification(_ proto.Request) {
}
