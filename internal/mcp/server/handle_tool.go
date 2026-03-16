//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/config/mcp/mime"
	"github.com/ActiveMemory/ctx/internal/config/mcp/tool"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// handleToolsList returns all available MCP tools.
func (s *Server) handleToolsList(req proto.Request) *proto.Response {
	return s.ok(req.ID, proto.ToolListResult{Tools: proto.ToolDefs})
}

// handleToolsCall dispatches a tool call to the appropriate handler.
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
func (s *Server) toolResult(
	id json.RawMessage, text string, err error,
) *proto.Response {
	if err != nil {
		return s.toolError(id, err.Error())
	}
	return s.toolOK(id, text)
}

// call invokes a no-arg handler and wraps the result.
func (s *Server) call(
	id json.RawMessage, fn func() (string, error),
) *proto.Response {
	text, err := fn()
	return s.toolResult(id, text, err)
}

// toolAdd extracts MCP args and delegates to handler.Add.
func (s *Server) toolAdd(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	entryType, content, errResp := s.extractEntryArgs(id, args)
	if errResp != nil {
		return errResp
	}
	text, err := s.handler.Add(entryType, content, extractOpts(args))
	return s.toolResult(id, text, err)
}

// toolComplete extracts the query and delegates to handler.Complete.
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
func (s *Server) toolRecall(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	limit := cfg.DefaultRecallLimit
	if v, ok := args[field.Limit].(float64); ok && v > 0 {
		limit = int(v)
	}

	sinceStr, _ := args[field.Since].(string)
	since, parseErr := handler.ParseRecallSince(sinceStr)
	if parseErr != nil {
		return s.toolError(
			id, fmt.Sprintf(
				assets.TextDesc(assets.TextDescKeyMCPInvalidSinceDate),
				parseErr,
			),
		)
	}

	text, err := s.handler.Recall(limit, since)
	return s.toolResult(id, text, err)
}

// toolWatchUpdate extracts MCP args and delegates to
// handler.WatchUpdate.
func (s *Server) toolWatchUpdate(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	entryType, content, errResp := s.extractEntryArgs(id, args)
	if errResp != nil {
		return errResp
	}
	text, err := s.handler.WatchUpdate(
		entryType, content, extractOpts(args),
	)
	return s.toolResult(id, text, err)
}

// toolCompact extracts the archive flag and delegates to
// handler.Compact.
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
func (s *Server) toolCheckTaskCompletion(
	id json.RawMessage, args map[string]interface{},
) *proto.Response {
	recentAction, _ := args[field.RecentAction].(string)
	text, err := s.handler.CheckTaskCompletion(recentAction)
	return s.toolResult(id, text, err)
}

// toolSessionEvent extracts event type/caller and delegates to
// handler.SessionEvent.
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

// extractEntryArgs validates required type/content from MCP args.
func (s *Server) extractEntryArgs(
	id json.RawMessage, args map[string]interface{},
) (entryType, content string, errResp *proto.Response) {
	entryType, _ = args[cli.AttrType].(string)
	content, _ = args[field.Content].(string)

	if entryType == "" || content == "" {
		return "", "", s.toolError(
			id, assets.TextDesc(assets.TextDescKeyMCPTypeContentRequired),
		)
	}

	return entryType, content, nil
}

// extractOpts builds EntryOpts from MCP tool arguments.
func extractOpts(args map[string]interface{}) handler.EntryOpts {
	opts := handler.EntryOpts{}
	if v, ok := args[field.Priority].(string); ok {
		opts.Priority = v
	}
	if v, ok := args[cli.AttrContext].(string); ok {
		opts.Context = v
	}
	if v, ok := args[cli.AttrRationale].(string); ok {
		opts.Rationale = v
	}
	if v, ok := args[cli.AttrConsequences].(string); ok {
		opts.Consequences = v
	}
	if v, ok := args[cli.AttrLesson].(string); ok {
		opts.Lesson = v
	}
	if v, ok := args[cli.AttrApplication].(string); ok {
		opts.Application = v
	}
	return opts
}

// toolOK builds a successful tool result.
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
func (s *Server) toolError(id json.RawMessage, msg string) *proto.Response {
	return s.ok(id, proto.CallToolResult{
		Content: []proto.ToolContent{{Type: mime.ContentTypeText, Text: msg}},
		IsError: true,
	})
}
