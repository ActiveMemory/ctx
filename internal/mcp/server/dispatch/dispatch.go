//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dispatch

import (
	"github.com/ActiveMemory/ctx/internal/config/mcp/method"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/dispatch/poll"
	"github.com/ActiveMemory/ctx/internal/mcp/server/ping"
	"github.com/ActiveMemory/ctx/internal/mcp/server/resource"
	"github.com/ActiveMemory/ctx/internal/mcp/server/route/fallback"
	"github.com/ActiveMemory/ctx/internal/mcp/server/route/initialize"
	"github.com/ActiveMemory/ctx/internal/mcp/server/route/prompt"
	"github.com/ActiveMemory/ctx/internal/mcp/server/route/tool"
)

// Do routes a request to the correct handler based on the
// method name.
//
// Parameters:
//   - version: server version string
//   - d: runtime dependencies for domain logic (context dir, budget, session)
//   - resList: pre-built resource list
//   - poller: resource poller for subscribe/unsubscribe
//   - req: parsed JSON-RPC request
//
// Returns:
//   - *proto.Response: result or error response
func Do(
	version string, d *entity.MCPDeps,
	resList proto.ResourceListResult, poller *poll.Poller,
	req proto.Request,
) *proto.Response {
	switch req.Method {
	case method.Initialize:
		return initialize.Dispatch(version, req)
	case method.Ping:
		return ping.Dispatch(req)
	case method.ResourceList:
		return resource.DispatchList(req, resList)
	case method.ResourceRead:
		return resource.DispatchRead(
			d.ContextDir, d.TokenBudget, req,
		)
	case method.ResourceSubscribe:
		return resource.DispatchSubscribe(req, poller.Subscribe)
	case method.ResourceUnsubscribe:
		return resource.DispatchUnsubscribe(req, poller.Unsubscribe)
	case method.ToolList:
		return tool.DispatchList(req)
	case method.ToolCall:
		return tool.DispatchCall(d, req)
	case method.PromptList:
		return prompt.DispatchList(req)
	case method.PromptGet:
		return prompt.DispatchGet(d, req)
	default:
		return fallback.DispatchErr(req)
	}
}
