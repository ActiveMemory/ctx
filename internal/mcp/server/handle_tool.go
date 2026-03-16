//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	timeCfg "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/mcp/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/extract"
)

// handleToolsList returns all available MCP tools.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *proto.Response: tool list result
func (s *Server) handleToolsList(req proto.Request) *proto.Response {
	return s.ok(req.ID, proto.ToolListResult{Tools: entity.ToolDefs})
}

// handleToolsCall dispatches a tool call to the appropriate handler.
//
// Parameters:
//   - req: the MCP request containing tool name and arguments
//
// Returns:
//   - *proto.Response: tool result or error
func (s *Server) handleToolsCall(req proto.Request) *proto.Response {
	var params proto.CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.error(
			req.ID, proto.ErrCodeInvalidArg,
			assets.TextDesc(assets.TextDescKeyMCPInvalidParams),
		)
	}

	s.handler.Session.RecordToolCall()

	switch params.Name {
	case tool.Status:
		return s.call(req.ID, s.handler.Status)
	case tool.Add:
		return s.toolAdd(req.ID, params.Arguments)
	case tool.Complete:
		return s.toolComplete(req.ID, params.Arguments)
	case tool.Drift:
		return s.call(req.ID, s.handler.Drift)
	case tool.Recall:
		return s.toolRecall(req.ID, params.Arguments)
	case tool.WatchUpdate:
		return s.toolWatchUpdate(req.ID, params.Arguments)
	case tool.Compact:
		return s.toolCompact(req.ID, params.Arguments)
	case tool.Next:
		return s.call(req.ID, s.handler.Next)
	case tool.CheckTaskCompletion:
		return s.toolCheckTaskCompletion(req.ID, params.Arguments)
	case tool.SessionEvent:
		return s.toolSessionEvent(req.ID, params.Arguments)
	case tool.Remind:
		return s.call(req.ID, s.handler.Remind)
	default:
		return s.error(
			req.ID, proto.ErrCodeNotFound,
			fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPUnknownTool),
				params.Name,
			),
		)
	}
}

// toolResult wraps a handler (string, error) return into a
// proto.Response.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - text: success text from the handler
//   - err: handler error, nil on success
//
// Returns:
//   - *proto.Response: tool OK or tool error response
func (s *Server) toolResult(
	id json.RawMessage, text string, err error,
) *proto.Response {
	if err != nil {
		return s.toolError(id, err.Error())
	}
	return s.toolOK(id, text)
}

// call invokes a no-arg handler and wraps the result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - fn: handler function returning (string, error)
//
// Returns:
//   - *proto.Response: wrapped handler result
func (s *Server) call(
	id json.RawMessage, fn func() (string, error),
) *proto.Response {
	text, err := fn()
	return s.toolResult(id, text, err)
}

// toolAdd extracts MCP args and delegates to handler.Add.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, content, optional fields)
//
// Returns:
//   - *proto.Response: add confirmation or validation error
func (s *Server) toolAdd(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	entryType, content, extractErr := extract.EntryArgs(args)
	if extractErr != nil {
		return s.toolError(id, extractErr.Error())
	}
	text, err := s.handler.Add(entryType, content, extract.Opts(args))
	return s.toolResult(id, text, err)
}

// toolComplete extracts the query and delegates to handler.Complete.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (query)
//
// Returns:
//   - *proto.Response: completion confirmation or error
func (s *Server) toolComplete(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	query, _ := args[field.Query].(string)
	if query == "" {
		return s.toolError(
			id, assets.TextDesc(assets.TextDescKeyMCPQueryRequired),
		)
	}
	text, err := s.handler.Complete(query)
	return s.toolResult(id, text, err)
}

// toolRecall extracts limit/since and delegates to handler.Recall.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (limit, since)
//
// Returns:
//   - *proto.Response: session list or parse error
func (s *Server) toolRecall(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	limit := cfg.DefaultRecallLimit
	if v, ok := args[field.Limit].(float64); ok && v > 0 {
		limit = int(v)
	}

	var since time.Time
	if sinceStr, _ := args[field.Since].(string); sinceStr != "" {
		var parseErr error
		since, parseErr = time.Parse(timeCfg.DateFormat, sinceStr)
		if parseErr != nil {
			return s.toolError(
				id, fmt.Sprintf(
					assets.TextDesc(assets.TextDescKeyMCPInvalidSinceDate),
					parseErr,
				),
			)
		}
	}

	text, err := s.handler.Recall(limit, since)
	return s.toolResult(id, text, err)
}

// toolWatchUpdate extracts MCP args and delegates to
// handler.WatchUpdate.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, content, optional fields)
//
// Returns:
//   - *proto.Response: write confirmation or validation error
func (s *Server) toolWatchUpdate(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	entryType, content, extractErr := extract.EntryArgs(args)
	if extractErr != nil {
		return s.toolError(id, extractErr.Error())
	}
	text, err := s.handler.WatchUpdate(
		entryType, content, extract.Opts(args),
	)
	return s.toolResult(id, text, err)
}

// toolCompact extracts the archive flag and delegates to
// handler.Compact.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (archive)
//
// Returns:
//   - *proto.Response: compact summary or error
func (s *Server) toolCompact(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	archive := false
	if v, ok := args[field.Archive].(bool); ok {
		archive = v
	}
	text, err := s.handler.Compact(archive)
	return s.toolResult(id, text, err)
}

// toolCheckTaskCompletion extracts recent_action and delegates to
// handler.CheckTaskCompletion.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (recent_action)
//
// Returns:
//   - *proto.Response: matching task prompt or empty result
func (s *Server) toolCheckTaskCompletion(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	recentAction, _ := args[field.RecentAction].(string)
	text, err := s.handler.CheckTaskCompletion(recentAction)
	return s.toolResult(id, text, err)
}

// toolSessionEvent extracts the event type/caller and delegates to
// handler.SessionEvent.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, caller)
//
// Returns:
//   - *proto.Response: session event confirmation or error
func (s *Server) toolSessionEvent(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	eventType, _ := args[cli.AttrType].(string)
	if eventType == "" {
		return s.toolError(id, assets.TextDesc(
			assets.TextDescKeyMCPEventTypeRequired),
		)
	}
	caller, _ := args[field.Caller].(string)
	text, err := s.handler.SessionEvent(eventType, caller)
	return s.toolResult(id, text, err)
}

// toolOK builds a successful tool result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - text: success text to include in the result
//
// Returns:
//   - *proto.Response: tool result with text content
func (s *Server) toolOK(id json.RawMessage, text string) *proto.Response {
	return s.ok(
		id,
		proto.CallToolResult{
			Content: []proto.ToolContent{
				{Type: mime.ContentTypeText, Text: text},
			},
		})
}

// toolError builds a tool error result.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - msg: error message text
//
// Returns:
//   - *proto.Response: tool result with IsError set
func (s *Server) toolError(id json.RawMessage, msg string) *proto.Response {
	return s.ok(id, proto.CallToolResult{
		Content: []proto.ToolContent{{Type: mime.ContentTypeText, Text: msg}},
		IsError: true,
	})
}
