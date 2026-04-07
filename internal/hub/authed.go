//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// authedCtx adds the bearer token to outgoing metadata.
func (c *Client) authedCtx(
	ctx context.Context,
) context.Context {
	if c.token == "" {
		return ctx
	}
	md := metadata.Pairs(
		"authorization", bearerPrefix+c.token,
	)
	return metadata.NewOutgoingContext(ctx, md)
}
