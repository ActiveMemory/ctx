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
)

// TestIntegration_PublishAndSync spins up a hub, registers
// two clients, publishes from one, and verifies the other
// receives the entry via sync.
func TestIntegration_PublishAndSync(t *testing.T) {
	srv, conn, adminTok := startTestServer(t)
	_ = srv

	// Register client A.
	regA := callRegister(t, conn, adminTok, "alpha")

	// Register client B.
	regB := callRegister(t, conn, adminTok, "beta")

	ctxA := authedCtx(regA.ClientToken)
	ctxB := authedCtx(regB.ClientToken)

	// Client A publishes an entry.
	pubResp := &PublishResponse{}
	pubErr := conn.Invoke(ctxA,
		"/ctx.hub.v1.CtxHub/Publish",
		&PublishRequest{
			Entries: []PublishEntry{
				{
					ID:        "e1",
					Type:      "decision",
					Content:   "Use gRPC for hub",
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

	// Client B syncs and should see the entry.
	stream, syncErr := conn.NewStream(ctxB,
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Sync",
	)
	if syncErr != nil {
		t.Fatalf("Sync stream: %v", syncErr)
	}
	if sendErr := stream.SendMsg(
		&SyncRequest{SinceSequence: 0},
	); sendErr != nil {
		t.Fatalf("Sync send: %v", sendErr)
	}
	if closeErr := stream.CloseSend(); closeErr != nil {
		t.Fatalf("Sync close: %v", closeErr)
	}

	msg := &EntryMsg{}
	if recvErr := stream.RecvMsg(msg); recvErr != nil {
		t.Fatalf("Sync recv: %v", recvErr)
	}
	if msg.Content != "Use gRPC for hub" {
		t.Errorf("want 'Use gRPC for hub', got %q",
			msg.Content)
	}
	if msg.Origin != "alpha" {
		t.Errorf("want origin 'alpha', got %q",
			msg.Origin)
	}
}

// TestIntegration_IncrementalSync verifies that sync with
// a non-zero since_sequence only returns new entries.
func TestIntegration_IncrementalSync(t *testing.T) {
	srv, conn, adminTok := startTestServer(t)
	_ = srv

	reg := callRegister(t, conn, adminTok, "proj")
	ctx := authedCtx(reg.ClientToken)

	// Publish two entries.
	pubResp := &PublishResponse{}
	pubErr := conn.Invoke(ctx,
		"/ctx.hub.v1.CtxHub/Publish",
		&PublishRequest{
			Entries: []PublishEntry{
				{
					ID: "a", Type: "learning",
					Content:   "First",
					Origin:    "proj",
					Timestamp: time.Now().Unix(),
				},
				{
					ID: "b", Type: "learning",
					Content:   "Second",
					Origin:    "proj",
					Timestamp: time.Now().Unix(),
				},
			},
		},
		pubResp,
	)
	if pubErr != nil {
		t.Fatal(pubErr)
	}

	// Sync since sequence 1 — should only get "Second".
	stream, syncErr := conn.NewStream(ctx,
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Sync",
	)
	if syncErr != nil {
		t.Fatal(syncErr)
	}
	if sendErr := stream.SendMsg(
		&SyncRequest{SinceSequence: 1},
	); sendErr != nil {
		t.Fatal(sendErr)
	}
	_ = stream.CloseSend()

	msg := &EntryMsg{}
	if recvErr := stream.RecvMsg(msg); recvErr != nil {
		t.Fatal(recvErr)
	}
	if msg.Content != "Second" {
		t.Errorf("want 'Second', got %q", msg.Content)
	}
}

// TestIntegration_TypeFilter verifies that sync with type
// filters only returns matching entries.
func TestIntegration_TypeFilter(t *testing.T) {
	srv, conn, adminTok := startTestServer(t)
	_ = srv

	reg := callRegister(t, conn, adminTok, "proj")
	ctx := authedCtx(reg.ClientToken)

	pubResp := &PublishResponse{}
	_ = conn.Invoke(ctx,
		"/ctx.hub.v1.CtxHub/Publish",
		&PublishRequest{
			Entries: []PublishEntry{
				{
					ID: "d1", Type: "decision",
					Content:   "Use Go",
					Origin:    "proj",
					Timestamp: time.Now().Unix(),
				},
				{
					ID: "l1", Type: "learning",
					Content:   "Avoid mocks",
					Origin:    "proj",
					Timestamp: time.Now().Unix(),
				},
			},
		},
		pubResp,
	)

	// Sync only learnings.
	stream, _ := conn.NewStream(ctx,
		&grpc.StreamDesc{ServerStreams: true},
		"/ctx.hub.v1.CtxHub/Sync",
	)
	_ = stream.SendMsg(&SyncRequest{
		Types:         []string{"learning"},
		SinceSequence: 0,
	})
	_ = stream.CloseSend()

	msg := &EntryMsg{}
	recvErr := stream.RecvMsg(msg)
	if recvErr != nil {
		t.Fatal(recvErr)
	}
	if msg.Type != "learning" {
		t.Errorf("want type 'learning', got %q", msg.Type)
	}
	if msg.Content != "Avoid mocks" {
		t.Errorf("want 'Avoid mocks', got %q",
			msg.Content)
	}
}

// TestIntegration_ClientLib verifies the Client library
// works end-to-end with the server.
func TestIntegration_ClientLib(t *testing.T) {
	_, _, adminTok := startTestServer(t)

	// Need a separate connection for the Client lib test
	// since startTestServer returns a raw conn.
	dir := t.TempDir()
	store, storeErr := NewStore(dir)
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	srv := NewServer(store, adminTok)
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(func() { srv.GracefulStop() })

	addr := lis.Addr().String()

	// Register via Client lib.
	client, dialErr := NewClient(addr, "")
	if dialErr != nil {
		t.Fatal(dialErr)
	}
	defer func() { _ = client.Close() }()

	regResp, regErr := client.Register(
		context.Background(), adminTok, "test-proj",
	)
	if regErr != nil {
		t.Fatalf("Register: %v", regErr)
	}

	// Reconnect with the client token.
	_ = client.Close()
	client2, dial2Err := NewClient(
		addr, regResp.ClientToken,
	)
	if dial2Err != nil {
		t.Fatal(dial2Err)
	}
	defer func() { _ = client2.Close() }()

	// Publish.
	_, pubErr := client2.Publish(
		context.Background(),
		[]PublishEntry{
			{
				ID: "c1", Type: "convention",
				Content:   "Use snake_case",
				Origin:    "test-proj",
				Timestamp: time.Now().Unix(),
			},
		},
	)
	if pubErr != nil {
		t.Fatalf("Publish: %v", pubErr)
	}

	// Sync.
	entries, syncErr := client2.Sync(
		context.Background(), nil, 0,
	)
	if syncErr != nil {
		t.Fatalf("Sync: %v", syncErr)
	}
	if len(entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(entries))
	}
	if entries[0].Content != "Use snake_case" {
		t.Errorf("want 'Use snake_case', got %q",
			entries[0].Content)
	}

	// Status.
	statusResp, statusErr := client2.Status(
		context.Background(),
	)
	if statusErr != nil {
		t.Fatalf("Status: %v", statusErr)
	}
	if statusResp.TotalEntries != 1 {
		t.Errorf("want 1 total, got %d",
			statusResp.TotalEntries)
	}
}

// startTestServer is defined in server_test.go — reused
// here for integration tests. The helper creates a
// temporary store, generates an admin token, starts the
// server on a random port, and returns a connected client.

// authedCtx and callRegister are also in server_test.go.
