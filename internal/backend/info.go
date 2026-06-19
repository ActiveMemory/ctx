//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backend

import "context"

// PingInfo checks backend reachability and returns optional model info.
//
// Parameters:
//   - ctx: request context
//   - target: backend to ping
//
// Returns:
//   - Response: optional ping metadata, including FirstModel when available
//   - error: backend ping failure
func PingInfo(ctx context.Context, target Backend) (Response, error) {
	if modelBackend, ok := target.(interface {
		models(context.Context) (Response, error)
	}); ok {
		return modelBackend.models(ctx)
	}
	return Response{}, target.Ping(ctx)
}

// EndpointInfo returns the resolved endpoint for HTTP-backed backends.
//
// Parameters:
//   - target: backend to inspect
//
// Returns:
//   - string: resolved endpoint when the backend exposes one
func EndpointInfo(target Backend) string {
	if endpointBackend, ok := target.(interface{ endpoint() string }); ok {
		return endpointBackend.endpoint()
	}
	return ""
}
