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

	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
)

// clientIDBytes is the size of a generated client UUID.
const clientIDBytes = 16

// generateClientID returns a hex-encoded random client ID.
func generateClientID() (string, error) {
	b := make([]byte, clientIDBytes)
	if _, randErr := rand.Read(b); randErr != nil {
		return "", errHub.GenerateToken(randErr)
	}
	return hex.EncodeToString(b), nil
}

// hubServiceName is the gRPC service descriptor name.
const hubServiceName = "ctx.hub.v1.CtxHub"

// serviceDesc returns the gRPC service description.
func serviceDesc(s *Server) *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: hubServiceName,
		HandlerType: (*any)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Register",
				Handler:    makeRegisterHandler(s),
			},
			{
				MethodName: "Publish",
				Handler:    makePublishHandler(s),
			},
			{
				MethodName: "Status",
				Handler:    makeStatusHandler(s),
			},
		},
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "Sync",
				Handler:       makeSyncHandler(s),
				ServerStreams: true,
			},
			{
				StreamName:    "Listen",
				Handler:       makeListenHandler(s),
				ServerStreams: true,
			},
		},
		Metadata: "hub.proto",
	}
}

// registerService registers the hub on a gRPC server.
func registerService(gs *grpc.Server, s *Server) {
	gs.RegisterService(serviceDesc(s), s)
}

// makeRegisterHandler creates the Register handler.
// Register uses admin token auth, not bearer.
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
