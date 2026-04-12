//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"google.golang.org/grpc"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
)

// generateClientID returns a hex-encoded random client ID.
//
// Returns:
//   - string: hex-encoded UUID
//   - error: non-nil if crypto/rand fails
func generateClientID() (string, error) {
	b := make([]byte, cfgHub.ClientIDBytes)
	if _, randErr := rand.Read(b); randErr != nil {
		return "", errHub.GenerateToken(randErr)
	}
	return hex.EncodeToString(b), nil
}

// serviceDesc returns the gRPC service description.
//
// Parameters:
//   - s: hub server to build descriptors for
//
// Returns:
//   - *grpc.ServiceDesc: gRPC service descriptor
func serviceDesc(s *Server) *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: cfgHub.ServiceName,
		HandlerType: (*any)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: cfgHub.MethodRegister,
				Handler:    makeRegisterHandler(s),
			},
			{
				MethodName: cfgHub.MethodPublish,
				Handler:    makePublishHandler(s),
			},
			{
				MethodName: cfgHub.MethodStatus,
				Handler:    makeStatusHandler(s),
			},
		},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    cfgHub.MethodSync,
				Handler:       makeSyncHandler(s),
				ServerStreams: true,
			},
			{
				StreamName:    cfgHub.MethodListen,
				Handler:       makeListenHandler(s),
				ServerStreams: true,
			},
		},
		Metadata: cfgHub.ProtoFile,
	}
}

// registerService registers the hub on a gRPC server.
//
// Parameters:
//   - gs: gRPC server to register on
//   - s: hub server providing RPC handlers
func registerService(gs *grpc.Server, s *Server) {
	gs.RegisterService(serviceDesc(s), s)
}

// makeRegisterHandler creates the Register handler.
// Register uses admin token auth, not bearer.
//
// Parameters:
//   - s: hub server for request dispatch
//
// Returns:
//   - grpc.MethodHandler: unary handler for Register RPC
func makeRegisterHandler(s *Server) grpc.MethodHandler {
	return func(
		_ any, ctx context.Context,
		dec func(any) error,
		_ grpc.UnaryServerInterceptor,
	) (any, error) {
		req := &RegisterRequest{}
		if decErr := dec(req); decErr != nil {
			return nil, decErr
		}
		return s.register(ctx, req)
	}
}

// makePublishHandler creates the Publish handler.
//
// Parameters:
//   - s: hub server for request dispatch
//
// Returns:
//   - grpc.MethodHandler: unary handler for Publish RPC
func makePublishHandler(s *Server) grpc.MethodHandler {
	return func(
		_ any, ctx context.Context,
		dec func(any) error,
		_ grpc.UnaryServerInterceptor,
	) (any, error) {
		if authErr := validateBearer(
			ctx, s.store,
		); authErr != nil {
			return nil, authErr
		}
		req := &PublishRequest{}
		if decErr := dec(req); decErr != nil {
			return nil, decErr
		}
		return s.publish(ctx, req)
	}
}

// makeStatusHandler creates the Status handler.
//
// Parameters:
//   - s: hub server for request dispatch
//
// Returns:
//   - grpc.MethodHandler: unary handler for Status RPC
func makeStatusHandler(s *Server) grpc.MethodHandler {
	return func(
		_ any, ctx context.Context,
		_ func(any) error,
		_ grpc.UnaryServerInterceptor,
	) (any, error) {
		if authErr := validateBearer(
			ctx, s.store,
		); authErr != nil {
			return nil, authErr
		}
		return s.hubStatus(ctx)
	}
}

// makeSyncHandler creates the Sync stream handler.
//
// Parameters:
//   - s: hub server for request dispatch
//
// Returns:
//   - func(any, grpc.ServerStream) error: stream handler
func makeSyncHandler(
	s *Server,
) func(any, grpc.ServerStream) error {
	return func(_ any, ss grpc.ServerStream) error {
		if authErr := validateBearer(
			ss.Context(), s.store,
		); authErr != nil {
			return authErr
		}
		req := &SyncRequest{}
		if recvErr := ss.RecvMsg(req); recvErr != nil {
			return recvErr
		}
		return s.syncEntries(
			req, func(m *EntryMsg) error {
				return ss.SendMsg(m)
			},
		)
	}
}

// makeListenHandler creates the Listen stream handler.
//
// Parameters:
//   - s: hub server for request dispatch
//
// Returns:
//   - func(any, grpc.ServerStream) error: stream handler
func makeListenHandler(
	s *Server,
) func(any, grpc.ServerStream) error {
	return func(_ any, ss grpc.ServerStream) error {
		if authErr := validateBearer(
			ss.Context(), s.store,
		); authErr != nil {
			return authErr
		}
		req := &ListenRequest{}
		if recvErr := ss.RecvMsg(req); recvErr != nil {
			return recvErr
		}
		return s.listenEntries(
			req, func(m *EntryMsg) error {
				return ss.SendMsg(m)
			}, ss.Context(),
		)
	}
}
