//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"encoding/json"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/mcp/cfg"
	cfgSchema "github.com/ActiveMemory/ctx/internal/config/mcp/schema"
	"github.com/ActiveMemory/ctx/internal/context/load"
	"github.com/ActiveMemory/ctx/internal/mcp/proto"
	"github.com/ActiveMemory/ctx/internal/mcp/server/catalog"
	"github.com/ActiveMemory/ctx/internal/mcp/server/out"
	"github.com/ActiveMemory/ctx/internal/sanitize"
)

// DispatchList returns the pre-built resource list.
//
// Parameters:
//   - req: the MCP request
//   - list: pre-built resource list result
//
// Returns:
//   - *proto.Response: resource list response
func DispatchList(
	req proto.Request, list proto.ResourceListResult,
) *proto.Response {
	return out.OkResponse(req.ID, list)
}

// DispatchRead loads context and returns the requested resource
// content.
//
// Parameters:
//   - contextDir: path to the .context/ directory
//   - tokenBudget: token budget for agent packet assembly
//   - req: the MCP request containing the resource URI
//
// Returns:
//   - *proto.Response: resource content or error
func DispatchRead(
	contextDir string, tokenBudget int, req proto.Request,
) *proto.Response {
	var params proto.ReadResourceParams
	if unmarshalErr := json.Unmarshal(
		req.Params, &params,
	); unmarshalErr != nil {
		return out.ErrResponse(
			req.ID, cfgSchema.ErrCodeInvalidArg,
			desc.Text(text.DescKeyMCPErrInvalidParams),
		)
	}

	ctx, loadErr := load.Do(contextDir)
	if loadErr != nil {
		return out.ErrResponse(req.ID, cfgSchema.ErrCodeInternal,
			fmt.Sprintf(
				desc.Text(text.DescKeyMCPLoadContext),
				loadErr,
			))
	}

	// Individual file resource.
	if fileName := catalog.FileForURI(params.URI); fileName != "" {
		return readContextFile(req.ID, ctx, fileName, params.URI)
	}

	// Assembled agent packet.
	if params.URI == catalog.AgentURI() {
		return readAgentPacket(req.ID, ctx, tokenBudget)
	}

	return out.ErrResponse(req.ID, cfgSchema.ErrCodeInvalidArg,
		fmt.Sprintf(
			desc.Text(text.DescKeyMCPErrUnknownResource),
			sanitize.Reflect(params.URI, cfg.MaxURILen),
		))
}

// DispatchSubscribe parses subscribe params and calls the provided
// subscribe function with the validated URI.
//
// Parameters:
//   - req: the MCP request containing the resource URI
//   - fn: subscribe function to call with the URI
//
// Returns:
//   - *proto.Response: empty success or validation error
func DispatchSubscribe(
	req proto.Request, fn func(string),
) *proto.Response {
	return applySubscription(req, fn)
}

// DispatchUnsubscribe parses unsubscribe params and calls the
// provided unsubscribe function with the validated URI.
//
// Parameters:
//   - req: the MCP request containing the resource URI
//   - fn: unsubscribe function to call with the URI
//
// Returns:
//   - *proto.Response: empty success or validation error
func DispatchUnsubscribe(
	req proto.Request, fn func(string),
) *proto.Response {
	return applySubscription(req, fn)
}
