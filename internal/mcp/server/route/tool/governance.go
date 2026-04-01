//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"github.com/ActiveMemory/ctx/internal/mcp/handler"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// appendGovernance appends governance advisory warnings to a tool
// response. It modifies the response in-place by appending warning
// text to the first content item.
//
// Parameters:
//   - resp: the MCP response to augment
//   - toolName: name of the tool that was called
//   - h: handler providing session governance state
func appendGovernance(
	resp *proto.Response, toolName string, h *handler.Handler,
) {
	warning := h.Session.CheckGovernance(toolName)
	if warning == "" {
		return
	}
	result, ok := resp.Result.(proto.CallToolResult)
	if !ok || len(result.Content) == 0 {
		return
	}
	result.Content[0].Text += warning
	resp.Result = result
}
