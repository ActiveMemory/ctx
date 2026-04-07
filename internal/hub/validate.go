//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// bearerPrefix is stripped from the authorization header.
const bearerPrefix = "Bearer "

// validateBearer extracts and validates the bearer token
// from gRPC metadata.
func validateBearer(
	ctx context.Context, store *Store,
) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(
			codes.Unauthenticated, "missing metadata",
		)
	}

	vals := md.Get("authorization")
	if len(vals) == 0 {
		return status.Error(
			codes.Unauthenticated, "missing token",
		)
	}

	token := strings.TrimPrefix(vals[0], bearerPrefix)
	if store.ValidateToken(token) == nil {
		return status.Error(
			codes.Unauthenticated, "invalid token",
		)
	}
	return nil
}
