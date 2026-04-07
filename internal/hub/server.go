//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"net"

	"google.golang.org/grpc"
)

// NewServer creates a hub server backed by the given store.
//
// Parameters:
//   - store: append-only storage backend
//   - adminToken: token required for Register RPC
//
// Returns:
//   - *Server: configured server (call Serve to start)
func NewServer(store *Store, adminToken string) *Server {
	s := &Server{
		store:      store,
		adminToken: adminToken,
		listeners:  newFanOut(),
	}

	gs := grpc.NewServer()
	registerService(gs, s)
	s.grpc = gs

	return s
}

// Serve starts the gRPC server on the given listener.
//
// Parameters:
//   - lis: network listener to accept connections on
//
// Returns:
//   - error: non-nil if the server fails to start
func (s *Server) Serve(lis net.Listener) error {
	return s.grpc.Serve(lis)
}

// GracefulStop stops the server gracefully.
func (s *Server) GracefulStop() {
	s.grpc.GracefulStop()
}
