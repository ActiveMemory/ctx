//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package mcp

import (
	"io"
	"sync"
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
	session     *sessionState
	poller      *resourcePoller
}
