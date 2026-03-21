//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fallback

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// DispatchErr returns a method-not-found error for unrecognized
// methods.
//
// Parameters:
//   - req: the MCP request with the unrecognized method
//
// Returns:
//   - *proto.Response: method-not-found error response
func DispatchErr(req proto.Request) *proto.Response {
	return out.ErrResponse(req.ID, proto.ErrCodeNotFound,
		fmt.Sprintf(
			desc.Text(text.DescKeyMCPErrMethodNotFound),
			req.Method,
		),
	)
}
