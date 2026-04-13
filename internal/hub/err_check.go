//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// authErr reports whether err is an authentication or
// authorization failure.
//
// Parameters:
//   - err: error to check
//
// Returns:
//   - bool: true if err is Unauthenticated or PermissionDenied
func authErr(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	c := s.Code()
	return c == codes.Unauthenticated ||
		c == codes.PermissionDenied
}
