//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tool

import (
	"github.com/ActiveMemory/ctx/internal/entity"
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
//   - d: runtime dependencies carrying the session state
func appendGovernance(
	resp *proto.Response, toolName string, d *entity.MCPDeps,
) {
	warning := handler.CheckGovernance(d, toolName)
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
