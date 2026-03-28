//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"io"
	"sync"

	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/poll"
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
	poller       *poll.Poller
	resourceList proto.ResourceListResult // pre-built, immutable after init
}
