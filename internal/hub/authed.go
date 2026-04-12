//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"google.golang.org/grpc/metadata"
)

// addBearerMD adds a bearer token to outgoing gRPC
// metadata.
//
// Parameters:
//   - ctx: parent context
//   - tok: bearer token to attach
//
// Returns:
//   - context.Context: context with bearer metadata added
func addBearerMD(
	ctx context.Context, tok string,
) context.Context {
	if tok == "" {
		return ctx
	}
	return metadata.NewOutgoingContext(
		ctx, metadata.Pairs(
			cfgHub.HeaderAuthorization, bearerPrefix+tok,
		),
	)
}

// authedCtx adds the bearer token to outgoing metadata.
//
// Parameters:
//   - ctx: parent context
//
// Returns:
//   - context.Context: context with bearer metadata added
func (c *Client) authedCtx(
	ctx context.Context,
) context.Context {
	if c.token == "" {
		return ctx
	}
	md := metadata.Pairs(
		cfgHub.HeaderAuthorization, bearerPrefix+c.token,
	)
	return metadata.NewOutgoingContext(ctx, md)
}
