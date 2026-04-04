//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/config/mcp/server"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
)

// Request unmarshals raw JSON into a proto.Request.
//
// Returns nil request for notifications (no ID). Returns an error
// response for malformed JSON.
//
// Parameters:
//   - data: raw JSON bytes from stdin
//
// Returns:
//   - *proto.Request: parsed request, nil for notifications
//   - *proto.Response: parse error response, nil on success
func Request(data []byte) (*proto.Request, *proto.Response) {
	var req proto.Request
	if unmarshalErr := json.Unmarshal(data, &req); unmarshalErr != nil {
		return nil, &proto.Response{
			JSONRPC: server.JSONRPCVersion,
			Error: &proto.RPCError{
				Code:    cfgSchema.ErrCodeParse,
				Message: desc.Text(text.DescKeyMCPErrParse),
			},
		}
	}

	// Notifications have no ID and expect no response.
	if req.ID == nil {
		return nil, nil
	}

	return &req, nil
}
