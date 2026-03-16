//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

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
	return out.OkResponse(req.ID, result)
}
