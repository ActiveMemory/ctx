//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/field"
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// steeringGet extracts the optional prompt and delegates to
// handler.SteeringGet.
//
// Parameters:
//   - h: handler for domain logic
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (prompt)
//
// Returns:
//   - *proto.Response: steering files or error
func steeringGet(
	h *handler.Handler, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	prompt, _ := args[field.Prompt].(string)
	t, err := h.SteeringGet(prompt)
	return out.ToolResult(id, t, err)
}

// search extracts the required query and delegates to
// handler.Search.
//
// Parameters:
//   - h: handler for domain logic
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (query)
//
// Returns:
//   - *proto.Response: search results or error
func search(
	h *handler.Handler, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	query, _ := args[field.Query].(string)
	if query == "" {
		return out.ToolError(
			id, desc.Text(text.DescKeyMCPErrQueryRequired),
		)
	}
	t, err := h.Search(query)
	return out.ToolResult(id, t, err)
}

// sessionEnd extracts the optional summary and delegates to
// handler.SessionEndHooks.
//
// Parameters:
//   - h: handler for domain logic
//   - id: JSON-RPC request ID
//   - args: MCP tool arguments (summary)
//
// Returns:
//   - *proto.Response: session end result or error
func sessionEnd(
	h *handler.Handler, id json.RawMessage,
	args map[string]interface{},
) *proto.Response {
	summary, _ := args[field.Summary].(string)
	t, err := h.SessionEndHooks(summary)
	return out.ToolResult(id, t, err)
}
