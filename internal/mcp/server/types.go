//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"io"

	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	mcpIO "github.com/ActiveMemory/ctx/internal/mcp/server/io"
	"github.com/ActiveMemory/ctx/internal/mcp/server/poll"
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
//   - handler: Domain logic handler for tool/resource/prompt calls
//   - version: Binary version for server info response
//   - out: Thread-safe JSON writer for stdout
//   - in: Input reader for stdin
//   - poller: Background resource change poller
//   - resourceList: Pre-built resource list (immutable after init)
type Server struct {
	handler      *handler.Handler
	version      string
	out          *mcpIO.Writer
	in           io.Reader
	poller       *poll.Poller
	resourceList proto.ResourceListResult // pre-built, immutable after init
}
