//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// Dispatch responds to the MCP initialize handshake.
//
// Parameters:
//   - version: server version string
//   - req: parsed JSON-RPC request
//
// Returns:
//   - *proto.Response: server capabilities and protocol version
func Dispatch(version string, req proto.Request) *proto.Response {
	return out.OkResponse(req.ID, proto.InitializeResult{
		ProtocolVersion: cfgSchema.ProtocolVersion,
		Capabilities: proto.ServerCaps{
			Resources: &proto.ResourcesCap{Subscribe: true},
			Tools:     &proto.ToolsCap{},
			Prompts:   &proto.PromptsCap{},
		},
		ServerInfo: proto.AppInfo{
			Name:    server.Name,
			Version: version,
		},
	})
}
