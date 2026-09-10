//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"errors"

	"github.com/hashicorp/raft"
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

// clusterOpErr maps a Raft configuration or leadership error
// onto a gRPC status.
//
// raft.ErrNotLeader is a precondition the operator can act on
// -- ask the leader instead -- so it travels as
// FailedPrecondition with an answer, rather than as an opaque
// internal error. Everything else is the hub's problem, not
// the caller's.
//
// Parameters:
//   - err: error from a Cluster configuration or transfer call
//
// Returns:
//   - error: gRPC status carrying the right code
func clusterOpErr(err error) error {
	if errors.Is(err, raft.ErrNotLeader) {
		return status.Error(
			codes.FailedPrecondition, cfgHub.ErrNotLeader,
		)
	}

	return status.Error(codes.Internal, err.Error())
}

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
