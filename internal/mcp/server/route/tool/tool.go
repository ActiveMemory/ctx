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
	"github.com/ActiveMemory/ctx/internal/config/entry"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/extract"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	"github.com/ActiveMemory/ctx/internal/sanitize"
)

// add extracts MCP args and delegates to [handler.Add].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, content, optional fields)
//
// Returns:
//   - *proto.Response: add confirmation or validation error
func add(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	entryType, content, extractErr := extract.EntryArgs(args)
	if extractErr != nil {
		return out.ToolError(id, extractErr.Error())
	}
	// MCP-SAN.2: Reject unknown entry types before writing.
	if _, ok := entry.CtxFile(entryType); !ok {
		return out.ToolError(id, fmt.Sprintf(
			desc.Text(text.DescKeyMCPErrUnknownEntryType),
			sanitize.Reflect(entryType, cfg.MaxNameLen),
		))
	}

	// MCP-SAN.3: Sanitize content before writing to .context/.
	content = sanitize.Content(content)

	t, addErr := handler.Add(
		d, entryType, content, extract.SanitizedOpts(args),
	)
	return out.ToolResult(id, t, addErr)
}

// complete extracts the query and delegates to [handler.Complete].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (query)
//
// Returns:
//   - *proto.Response: completion confirmation or error
func complete(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	query, _ := args[field.Query].(string)
	if query == "" {
		return out.ToolError(
			id, desc.Text(text.DescKeyMCPErrQueryRequired),
		)
	}
	query = sanitize.Reflect(query, cfg.MaxQueryLen)
	t, completeErr := handler.Complete(d, query)
	return out.ToolResult(id, t, completeErr)
}

// journalSource extracts limit/since and delegates to [handler.Recall].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (limit, since)
//
// Returns:
//   - *proto.Response: session list or parse error
func journalSource(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	limit := cfg.DefaultSourceLimit
	if v, ok := args[field.Limit].(float64); ok && v > 0 {
		limit = int(v)
	}

	// MCP-SAN.1: Cap source limit to a reasonable upper bound.
	if limit > cfg.MaxSourceLimit {
		limit = cfg.MaxSourceLimit
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

	t, recallErr := handler.Recall(d, limit, since)
	return out.ToolResult(id, t, recallErr)
}

// watchUpdate extracts MCP args and delegates to
// [handler.WatchUpdate].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, content, optional fields)
//
// Returns:
//   - *proto.Response: write confirmation or validation error
func watchUpdate(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	entryType, content, extractErr := extract.EntryArgs(args)
	if extractErr != nil {
		return out.ToolError(id, extractErr.Error())
	}
	// MCP-SAN.2: Reject unknown entry types (allow "complete" as
	// special case handled by handler.WatchUpdate).
	if entryType != entry.Complete {
		if _, ok := entry.CtxFile(entryType); !ok {
			return out.ToolError(id, fmt.Sprintf(
				desc.Text(
					text.DescKeyMCPErrUnknownEntryType,
				),
				sanitize.Reflect(
					entryType, cfg.MaxNameLen,
				),
			))
		}
	}

	// MCP-SAN.3: Sanitize content before writing to .context/.
	content = sanitize.Content(content)

	t, updateErr := handler.WatchUpdate(
		d, entryType, content,
		extract.SanitizedOpts(args),
	)
	return out.ToolResult(id, t, updateErr)
}

// compact extracts the archive flag and delegates to [handler.Compact].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (archive)
//
// Returns:
//   - *proto.Response: compact summary or error
func compact(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	doArchive := false
	if v, ok := args[field.Archive].(bool); ok {
		doArchive = v
	}
	t, compactErr := handler.Compact(d, doArchive)
	return out.ToolResult(id, t, compactErr)
}

// checkTaskCompletion extracts recent_action and delegates to
// [handler.CheckTaskCompletion].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (recent_action)
//
// Returns:
//   - *proto.Response: matching task prompt or empty result
func checkTaskCompletion(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	recentAction, _ := args[field.RecentAction].(string)
	t, checkErr := handler.CheckTaskCompletion(d, recentAction)
	return out.ToolResult(id, t, checkErr)
}

// sessionEvent extracts the event type/caller and delegates to
// [handler.SessionEvent].
//
// Parameters:
//   - d: runtime dependencies
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (type, caller)
//
// Returns:
//   - *proto.Response: session event confirmation or error
func sessionEvent(
	d *entity.MCPDeps, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	eventType, _ := args[cli.AttrType].(string)
	if eventType == "" {
		return out.ToolError(id, desc.Text(
			text.DescKeyMCPEventTypeRequired),
		)
	}
	caller, _ := args[field.Caller].(string)

	// MCP-SAN.4: Sanitize caller before reflecting in response.
	caller = sanitize.Reflect(caller, cfg.MaxCallerLen)

	t, eventErr := handler.SessionEvent(d, eventType, caller)
	return out.ToolResult(id, t, eventErr)
}
