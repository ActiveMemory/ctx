//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/mcp"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// NewServer creates a new MCP server for the given context directory.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - version: binary version string for the server info response
//
// Returns:
//   - *Server: a configured MCP server ready to serve
func NewServer(contextDir, version string) *Server {
	return &Server{
		contextDir:  contextDir,
		version:     version,
		tokenBudget: rc.TokenBudget(),
		out:         os.Stdout,
		in:          os.Stdin,
	}
}

// Serve starts the MCP server, reading from stdin and writing to stdout.
//
// It blocks until stdin is closed or an unrecoverable error occurs.
// Each line from stdin is expected to be a JSON-RPC 2.0 request.
//
// Returns:
//   - error: non-nil if an I/O error prevents continued operation
func (s *Server) Serve() error {
	scanner := bufio.NewScanner(s.in)

	scanner.Buffer(make([]byte, 0, mcp.MCPScanMaxSize), mcp.MCPScanMaxSize)

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
			s.writeError(nil, errCodeInternal, assets.TextDesc(
				assets.TextDescKeyMCPFailedMarshal),
			)
			continue
		}
		if _, writeErr := s.out.Write(
			append(out, token.NewlineLF[0]),
		); writeErr != nil {
			return writeErr
		}
	}

	return scanner.Err()
}
