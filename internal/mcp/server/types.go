//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"io"

	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/dispatch/poll"
	mcpIO "github.com/ActiveMemory/ctx/internal/mcp/server/io"
)

// Server is an MCP server that exposes ctx context over JSON-RPC 2.0.
//
// It reads JSON-RPC requests from stdin and writes responses to stdout,
// following the Model Context Protocol specification.
//
// Thread-safety: out is a [mcpIO.Writer] that serializes all writes
// (main loop and poller goroutine). The main loop itself is
// single-threaded, so request dispatch and session mutations need
// no additional locking.
//
// Fields:
//   - deps: Runtime dependencies passed to every handler function
//     (context dir, token budget, session state)
//   - version: Binary version for server info response
//   - out: Thread-safe JSON writer for stdout
//   - in: Input reader for stdin
//   - poller: Background resource change poller
//   - resourceList: Pre-built resource list (immutable after init)
type Server struct {
	deps         *entity.MCPDeps
	version      string
	out          *mcpIO.Writer
	in           io.Reader
	poller       *poll.Poller
	resourceList proto.ResourceListResult // pre-built, immutable after init
}
