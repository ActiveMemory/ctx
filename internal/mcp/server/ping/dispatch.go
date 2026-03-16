//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ping

import (
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// Dispatch responds to a ping request with an empty success.
//
// Parameters:
//   - req: the MCP request
//
// Returns:
//   - *proto.Response: empty success response
func Dispatch(req proto.Request) *proto.Response {
	return out.OkResponse(req.ID, struct{}{})
}
