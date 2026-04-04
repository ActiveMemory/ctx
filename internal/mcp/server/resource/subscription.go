//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
)

// applySubscription handles the shared parse-validate-apply logic
// for subscribe and unsubscribe requests.
func applySubscription(
	req proto.Request, fn func(string),
) *proto.Response {
	var params proto.SubscribeParams
	if unmarshalErr := json.Unmarshal(
		req.Params, &params,
	); unmarshalErr != nil {
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams),
		)
	}
	if params.URI == "" {
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrURIRequired),
		)
	}
	fn(params.URI)
	return out.OkResponse(req.ID, struct{}{})
}
