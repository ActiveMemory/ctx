//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// startTestServer spins up a hub server on a random port.
func startTestServer(
	t *testing.T,
) (*Server, *grpc.ClientConn, string) {
	t.Helper()

	dir := t.TempDir()
	store, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	adminTok, err := GenerateAdminToken()
	if err != nil {
		t.Fatal(err)
	}

	srv := NewServer(store, adminTok)
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}

	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(func() { srv.GracefulStop() })

	conn, dialErr := grpc.NewClient(
		lis.Addr().String(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithDefaultCallOptions(
			grpc.CallContentSubtype(codecName),
		),
	)
	if dialErr != nil {
		t.Fatal(dialErr)
	}
	t.Cleanup(func() { _ = conn.Close() })

	return srv, conn, adminTok
}

// callRegister invokes the Register RPC directly.
func callRegister(
	t *testing.T,
	conn *grpc.ClientConn,
	adminTok, project string,
) *RegisterResponse {
	t.Helper()

	resp := &RegisterResponse{}
	err := conn.Invoke(
		context.Background(),
		"/ctx.hub.v1.CtxHub/Register",
		&RegisterRequest{
			AdminToken:  adminTok,
			ProjectName: project,
		},
		resp,
	)
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	return resp
}

// authedCtx returns a context with a bearer token.
func authedCtx(token string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+token)
	return metadata.NewOutgoingContext(
		context.Background(), md,
	)
}

func TestServerRegister(t *testing.T) {
	_, conn, adminTok := startTestServer(t)

	resp := callRegister(t, conn, adminTok, "alpha")
	if resp.ClientID == "" {
		t.Error("expected client ID")
	}
	if resp.ClientToken == "" {
		t.Error("expected client token")
	}
}

func TestServerRegisterBadToken(t *testing.T) {
	_, conn, _ := startTestServer(t)

	resp := &RegisterResponse{}
	err := conn.Invoke(
		context.Background(),
		"/ctx.hub.v1.CtxHub/Register",
		&RegisterRequest{
			AdminToken:  "wrong",
			ProjectName: "alpha",
		},
		resp,
	)
	if err == nil {
		t.Fatal("expected error for bad admin token")
	}
}

func TestServerPublishAndSync(t *testing.T) {
	_, conn, adminTok := startTestServer(t)

	reg := callRegister(t, conn, adminTok, "alpha")
	ctx := authedCtx(reg.ClientToken)

	// Publish
	pubResp := &PublishResponse{}
	pubErr := conn.Invoke(ctx,
		"/ctx.hub.v1.CtxHub/Publish",
		&PublishRequest{
			Entries: []PublishEntry{
				{
					ID:        "e1",
					Type:      "decision",
					Content:   "Use Go",
					Origin:    "alpha",
					Timestamp: time.Now().Unix(),
				},
			},
		},
		pubResp,
	)
	if pubErr != nil {
		t.Fatalf("Publish: %v", pubErr)
	}
	if len(pubResp.Sequences) != 1 {
		t.Fatalf("expected 1 sequence, got %d",
			len(pubResp.Sequences))
	}

	// Sync
	stream, syncErr := conn.NewStream(ctx,
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Sync",
	)
	if syncErr != nil {
		t.Fatalf("Sync stream: %v", syncErr)
	}
	if sendErr := stream.SendMsg(&SyncRequest{
		SinceSequence: 0,
	}); sendErr != nil {
		t.Fatalf("Sync send: %v", sendErr)
	}
	if closeErr := stream.CloseSend(); closeErr != nil {
		t.Fatalf("Sync close: %v", closeErr)
	}

	msg := &EntryMsg{}
	if recvErr := stream.RecvMsg(msg); recvErr != nil {
		t.Fatalf("Sync recv: %v", recvErr)
	}
	if msg.Content != "Use Go" {
		t.Errorf("expected 'Use Go', got %q", msg.Content)
	}
}

func TestServerStatus(t *testing.T) {
	_, conn, adminTok := startTestServer(t)

	reg := callRegister(t, conn, adminTok, "beta")
	ctx := authedCtx(reg.ClientToken)

	resp := &StatusResponse{}
	err := conn.Invoke(ctx,
		"/ctx.hub.v1.CtxHub/Status",
		&struct{}{},
		resp,
	)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if resp.TotalEntries != 0 {
		t.Errorf("expected 0 entries, got %d",
			resp.TotalEntries)
	}
}

func TestServerUnauthenticated(t *testing.T) {
	_, conn, adminTok := startTestServer(t)

	// Register to confirm the server works.
	callRegister(t, conn, adminTok, "test")

	// Publish with no auth should fail.
	resp := &PublishResponse{}
	err := conn.Invoke(
		context.Background(),
		"/ctx.hub.v1.CtxHub/Publish",
		&PublishRequest{},
		resp,
	)
	if err == nil {
		t.Fatal("expected unauthenticated error")
	}
	t.Logf("got expected error: %v", err)
}
