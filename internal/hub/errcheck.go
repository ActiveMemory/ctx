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

// isAuthErr reports whether err is an authentication or
// authorization failure.
func isAuthErr(err error) bool {
	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	c := s.Code()
	return c == codes.Unauthenticated ||
		c == codes.PermissionDenied
}
