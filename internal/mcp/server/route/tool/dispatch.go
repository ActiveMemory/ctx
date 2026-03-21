//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	tooldef "github.com/ActiveMemory/ctx/internal/mcp/server/def/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// DispatchList returns all available tools.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *proto.Response: tool list response
func DispatchList(req proto.Request) *proto.Response {
	return out.OkResponse(req.ID, proto.ToolListResult{Tools: tooldef.Defs})
}

// DispatchCall unmarshals tool call params and dispatches to the
// appropriate handler function.
//
// Parameters:
//   - h: handler for domain logic and session tracking
//   - req: the MCP request containing tool name and arguments
//
// Returns:
//   - *proto.Response: tool result or error
func DispatchCall(
	h *handler.Handler, req proto.Request,
) *proto.Response {
	var params proto.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return out.ErrResponse(
			req.ID, proto.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams),
		)
	}

	h.Session.RecordToolCall()

	switch params.Name {
	case tool.Status:
		return out.Call(req.ID, h.Status)
	case tool.Add:
		return add(h, req.ID, params.Arguments)
	case tool.Complete:
		return complete(h, req.ID, params.Arguments)
	case tool.Drift:
		return out.Call(req.ID, h.Drift)
	case tool.Recall:
		return recall(req.ID, params.Arguments, h.Recall)
	case tool.WatchUpdate:
		return watchUpdate(h, req.ID, params.Arguments)
	case tool.Compact:
		return compact(req.ID, params.Arguments, h.Compact)
	case tool.Next:
		return out.Call(req.ID, h.Next)
	case tool.CheckTaskCompletion:
		return checkTaskCompletion(
			req.ID, params.Arguments, h.CheckTaskCompletion,
		)
	case tool.SessionEvent:
		return sessionEvent(req.ID, params.Arguments, h.SessionEvent)
	case tool.Remind:
		return out.Call(req.ID, h.Remind)
	default:
		return out.ErrResponse(
			req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPErrUnknownTool),
				params.Name,
			),
		)
	}
}
