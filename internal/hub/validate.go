//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"strings"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// bearerPrefix is stripped from the authorization header.
const bearerPrefix = cfgHub.BearerPrefix

// validateBearer extracts and validates the bearer token
// from gRPC metadata.
//
// Parameters:
//   - ctx: request context with gRPC metadata
//   - store: store for token validation
//
// Returns:
//   - error: non-nil if token is missing or invalid
func validateBearer(
	ctx context.Context, store *Store,
) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(
			codes.Unauthenticated, cfgHub.ErrMissingMetadata,
		)
	}

	vals := md.Get(cfgHub.HeaderAuthorization)
	if len(vals) == 0 {
		return status.Error(
			codes.Unauthenticated, cfgHub.ErrMissingToken,
		)
	}

	token := strings.TrimPrefix(vals[0], bearerPrefix)
	if store.ValidateToken(token) == nil {
		return status.Error(
			codes.Unauthenticated, cfgHub.ErrInvalidToken,
		)
	}
	return nil
}
