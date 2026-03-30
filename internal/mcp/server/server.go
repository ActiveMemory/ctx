//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"bufio"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/catalog"
	"github.com/ActiveMemory/ctx/internal/mcp/server/dispatch"
	mcpIO "github.com/ActiveMemory/ctx/internal/mcp/server/io"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	"github.com/ActiveMemory/ctx/internal/mcp/server/parse"
	"github.com/ActiveMemory/ctx/internal/mcp/server/poll"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// New creates a new MCP server for the given context directory.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - version: binary version string for the server info response
//
// Returns:
//   - *Server: a configured MCP server ready to serve
func New(contextDir, version string) *Server {
	catalog.Init()
	srv := &Server{
		handler:      handler.New(contextDir, rc.TokenBudget()),
		version:      version,
		out:          mcpIO.NewWriter(os.Stdout),
		in:           os.Stdin,
		resourceList: catalog.ToList(),
	}
	srv.poller = poll.NewPoller(contextDir, func(n proto.Notification) {
		_ = srv.out.WriteJSON(n)
	})
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

		req, errResp := parse.Request(line)
		if errResp != nil {
			if writeErr := s.out.WriteJSON(errResp); writeErr != nil {
				return writeErr
			}
			continue
		}
		if req == nil {
			// Notification: no response required.
			continue
		}

		resp := dispatch.Do(
			s.version, s.handler, s.resourceList, s.poller, *req,
		)

		if writeErr := s.out.WriteJSON(resp); writeErr != nil {
			// Marshal failure: try to report it as an error response.
			fallback := out.ErrResponse(
				nil, proto.ErrCodeInternal,
				desc.Text(text.DescKeyMCPErrFailedMarshal),
			)
			if fbErr := s.out.WriteJSON(fallback); fbErr != nil {
				return fbErr
			}
			continue
		}
	}

	return scanner.Err()
}
