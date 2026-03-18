//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dispatch

import (
	"github.com/ActiveMemory/ctx/internal/config/mcp/method"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/ping"
	"github.com/ActiveMemory/ctx/internal/mcp/server/poll"
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
//   - h: handler for domain logic
//   - resList: pre-built resource list
//   - poller: resource poller for subscribe/unsubscribe
//   - req: parsed JSON-RPC request
//
// Returns:
//   - *proto.Response: result or error response
func Do(
	version string, h *handler.Handler,
	resList proto.ResourceListResult, poller *poll.ResourcePoller,
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
			h.ContextDir, h.TokenBudget, req,
		)
	case method.ResourceSubscribe:
		return resource.DispatchSubscribe(req, poller.Subscribe)
	case method.ResourceUnsubscribe:
		return resource.DispatchUnsubscribe(req, poller.Unsubscribe)
	case method.ToolList:
		return tool.DispatchList(req)
	case method.ToolCall:
		return tool.DispatchCall(h, req)
	case method.PromptList:
		return prompt.DispatchList(req)
	case method.PromptGet:
		return prompt.DispatchGet(h, req)
	default:
		return fallback.DispatchErr(req)
	}
}
