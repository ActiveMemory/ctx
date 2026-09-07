//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
)

// errSlowListener terminates a Listen stream whose fan-out
// channel [fanOut.broadcast] closed for being too slow. It is a
// package-level sentinel so [Server.listenEntries] returns the
// same value every time and tests can match it with errors.Is,
// while the ResourceExhausted code travels to the client: the
// stream ends with a reason instead of a silent EOF.
var errSlowListener = status.Error(
	codes.ResourceExhausted, cfgHub.ErrSlowListener,
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
