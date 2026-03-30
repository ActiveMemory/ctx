//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/cli"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/extract"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// add extracts MCP args and delegates to handler.Add.
//
// Parameters:
//   - h: handler for domain logic
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, content, optional fields)
//
// Returns:
//   - *proto.Response: add confirmation or validation error
func add(
	h *handler.Handler, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	entryType, content, extractErr := extract.EntryArgs(args)
	if extractErr != nil {
		return out.ToolError(id, extractErr.Error())
	}
	t, addErr := h.Add(entryType, content, extract.Opts(args))
	return out.ToolResult(id, t, addErr)
}

// complete extracts the query and delegates to handler.Complete.
//
// Parameters:
//   - h: handler for domain logic
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (query)
//
// Returns:
//   - *proto.Response: completion confirmation or error
func complete(
	h *handler.Handler, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	query, _ := args[field.Query].(string)
	if query == "" {
		return out.ToolError(
			id, desc.Text(text.DescKeyMCPErrQueryRequired),
		)
	}
	t, completeErr := h.Complete(query)
	return out.ToolResult(id, t, completeErr)
}

// journalSource extracts limit/since and calls the session query function.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (limit, since)
//   - fn: session query function accepting limit and since
//
// Returns:
//   - *proto.Response: session list or parse error
func journalSource(
	id json.RawMessage, args map[string]interface{},
	fn func(int, time.Time) (string, error),
) *proto.Response {
	limit := cfg.DefaultSourceLimit
	if v, ok := args[field.Limit].(float64); ok && v > 0 {
		limit = int(v)
	}

	var since time.Time
	if sinceStr, _ := args[field.Since].(string); sinceStr != "" {
		var parseErr error
		since, parseErr = time.Parse(cfgTime.DateFormat, sinceStr)
		if parseErr != nil {
			return out.ToolError(
				id, fmt.Sprintf(
					desc.Text(text.DescKeyMCPInvalidSinceDate),
					parseErr,
				),
			)
		}
	}

	t, recallErr := fn(limit, since)
	return out.ToolResult(id, t, recallErr)
}

// watchUpdate extracts MCP args and delegates to
// handler.WatchUpdate.
//
// Parameters:
//   - h: handler for domain logic
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, content, optional fields)
//
// Returns:
//   - *proto.Response: write confirmation or validation error
func watchUpdate(
	h *handler.Handler, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	entryType, content, extractErr := extract.EntryArgs(args)
	if extractErr != nil {
		return out.ToolError(id, extractErr.Error())
	}
	t, updateErr := h.WatchUpdate(
		entryType, content, extract.Opts(args),
	)
	return out.ToolResult(id, t, updateErr)
}

// compact extracts the archive flag and calls the compact
// function.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (archive)
//   - fn: compact function accepting archive flag
//
// Returns:
//   - *proto.Response: compact summary or error
func compact(
	id json.RawMessage, args map[string]interface{},
	fn func(bool) (string, error),
) *proto.Response {
	doArchive := false
	if v, ok := args[field.Archive].(bool); ok {
		doArchive = v
	}
	t, compactErr := fn(doArchive)
	return out.ToolResult(id, t, compactErr)
}

// checkTaskCompletion extracts recent_action and calls the
// check function.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (recent_action)
//   - fn: check function accepting action description
//
// Returns:
//   - *proto.Response: matching task prompt or empty result
func checkTaskCompletion(
	id json.RawMessage, args map[string]interface{},
	fn func(string) (string, error),
) *proto.Response {
	recentAction, _ := args[field.RecentAction].(string)
	t, checkErr := fn(recentAction)
	return out.ToolResult(id, t, checkErr)
}

// sessionEvent extracts the event type/caller and calls the
// session event function.
//
// Parameters:
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, caller)
//   - fn: session event function accepting type and caller
//
// Returns:
//   - *proto.Response: session event confirmation or error
func sessionEvent(
	id json.RawMessage, args map[string]interface{},
	fn func(string, string) (string, error),
) *proto.Response {
	eventType, _ := args[cli.AttrType].(string)
	if eventType == "" {
		return out.ToolError(id, desc.Text(
			text.DescKeyMCPEventTypeRequired),
		)
	}
	caller, _ := args[field.Caller].(string)
	t, eventErr := fn(eventType, caller)
	return out.ToolResult(id, t, eventErr)
}
