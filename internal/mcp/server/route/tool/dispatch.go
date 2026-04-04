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
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	defTool "github.com/ActiveMemory/ctx/internal/mcp/server/def/tool"
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
	return out.OkResponse(req.ID, proto.ToolListResult{Tools: defTool.Defs()})
}

// DispatchCall unmarshals tool call params and dispatches to the
// appropriate handler function. After dispatch, per-tool governance
// state is recorded and advisory warnings are appended to the
// response text.
//
// Parameters:
//   - h: handler for domain logic and session tracking
//   - req: the MCP request containing tool name and arguments
//
// Returns:
//   - *proto.Response: tool result or error (with governance warnings)
func DispatchCall(
	h *handler.Handler, req proto.Request,
) *proto.Response {
	var params proto.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams),
		)
	}

	h.Session.RecordToolCall()
	h.Session.IncrementCallsSinceWrite()

	var resp *proto.Response

	switch params.Name {
	case tool.Status:
		resp = out.Call(req.ID, h.Status)
		h.Session.RecordContextLoaded()
	case tool.Add:
		resp = add(h, req.ID, params.Arguments)
		h.Session.RecordContextWrite()
	case tool.Complete:
		resp = complete(h, req.ID, params.Arguments)
		h.Session.RecordContextWrite()
	case tool.Drift:
		resp = out.Call(req.ID, h.Drift)
		h.Session.RecordDriftCheck()
	case tool.JournalSource:
		resp = journalSource(req.ID, params.Arguments, h.Recall)
	case tool.WatchUpdate:
		resp = watchUpdate(h, req.ID, params.Arguments)
		h.Session.RecordContextWrite()
	case tool.Compact:
		resp = compact(req.ID, params.Arguments, h.Compact)
		h.Session.RecordContextWrite()
	case tool.Next:
		resp = out.Call(req.ID, h.Next)
	case tool.CheckTaskCompletion:
		resp = checkTaskCompletion(
			req.ID, params.Arguments, h.CheckTaskCompletion,
		)
	case tool.SessionEvent:
		resp = sessionEvent(req.ID, params.Arguments, h.SessionEvent)
	case tool.Remind:
		resp = out.Call(req.ID, h.Remind)
	case tool.SteeringGet:
		resp = steeringGet(h, req.ID, params.Arguments)
	case tool.Search:
		resp = search(h, req.ID, params.Arguments)
	case tool.SessionStart:
		resp = out.Call(req.ID, h.SessionStartHooks)
	case tool.SessionEnd:
		resp = sessionEnd(h, req.ID, params.Arguments)
	default:
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeNotFound,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPErrUnknownTool),
				params.Name,
			),
		)
	}

	appendGovernance(resp, params.Name, h)

	return resp
}
