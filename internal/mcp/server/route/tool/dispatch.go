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
	"github.com/ActiveMemory/ctx/internal/entity"
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
//   - d: runtime dependencies for domain logic and session tracking
//   - req: the MCP request containing tool name and arguments
//
// Returns:
//   - *proto.Response: tool result or error (with governance warnings)
func DispatchCall(
	d *entity.MCPDeps, req proto.Request,
) *proto.Response {
	var params proto.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams),
		)
	}

	d.Session.RecordToolCall()
	d.Session.IncrementCallsSinceWrite()

	var resp *proto.Response

	switch params.Name {
	case tool.Status:
		resp = out.Call(req.ID, func() (string, error) {
			return handler.Status(d)
		})
		d.Session.RecordContextLoaded()
	case tool.Add:
		resp = add(d, req.ID, params.Arguments)
		d.Session.RecordContextWrite()
	case tool.Complete:
		resp = complete(d, req.ID, params.Arguments)
		d.Session.RecordContextWrite()
	case tool.Drift:
		resp = out.Call(req.ID, func() (string, error) {
			return handler.Drift(d)
		})
		d.Session.RecordDriftCheck()
	case tool.JournalSource:
		resp = journalSource(d, req.ID, params.Arguments)
	case tool.WatchUpdate:
		resp = watchUpdate(d, req.ID, params.Arguments)
		d.Session.RecordContextWrite()
	case tool.Compact:
		resp = compact(d, req.ID, params.Arguments)
		d.Session.RecordContextWrite()
	case tool.Next:
		resp = out.Call(req.ID, func() (string, error) {
			return handler.Next(d)
		})
	case tool.CheckTaskCompletion:
		resp = checkTaskCompletion(d, req.ID, params.Arguments)
	case tool.SessionEvent:
		resp = sessionEvent(d, req.ID, params.Arguments)
	case tool.Remind:
		resp = out.Call(req.ID, func() (string, error) {
			return handler.Remind(d)
		})
	case tool.SteeringGet:
		resp = steeringGet(d, req.ID, params.Arguments)
	case tool.Search:
		resp = search(d, req.ID, params.Arguments)
	case tool.SessionStart:
		resp = out.Call(req.ID, func() (string, error) {
			return handler.SessionStartHooks(d)
		})
	case tool.SessionEnd:
		resp = sessionEnd(d, req.ID, params.Arguments)
	default:
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeNotFound,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPErrUnknownTool),
				params.Name,
			),
		)
	}

	appendGovernance(resp, params.Name, d)

	return resp
}
