//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
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

// slowListenerPayload is large enough that a handful of entries
// saturate the gRPC stream window once the client stops reading,
// which is what stalls the server's send and lets the fan-out
// buffer fill.
const slowListenerPayload = 128 << 10

// slowListenerBroadcastCap bounds the publish loop in
// [TestIntegration_SlowListenerReachesClient] so a regression
// fails the test instead of broadcasting forever.
const slowListenerBroadcastCap = 512

// TestIntegration_SlowListenerReachesClient is the client-visible
// half of the disconnect contract, over a real gRPC stream. A
// listener that stops draining is cut loose server-side; the
// client must learn that from a ResourceExhausted error, not from
// a stream that stays open forever while entries pass it by. This
// is what makes `ctx connection listen` exit non-zero rather than
// report success on a stream it is no longer receiving.
func TestIntegration_SlowListenerReachesClient(t *testing.T) {
	restore := logWarn.SetSink(io.Discard)
	defer restore()

	_, _, adminTok := startTestServer(t)

	store, storeErr := NewStore(t.TempDir())
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
	admin, adminDialErr := NewClient(addr, "")
	if adminDialErr != nil {
		t.Fatal(adminDialErr)
	}
	regResp, regErr := admin.Register(
		context.Background(), adminTok, "slow-proj",
	)
	if closeErr := admin.Close(); closeErr != nil {
		t.Log(closeErr)
	}
	if regErr != nil {
		t.Fatal(regErr)
	}

	client, dialErr := NewClient(addr, regResp.ClientToken)
	if dialErr != nil {
		t.Fatal(dialErr)
	}
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// The handler stalls on release, so the client stops calling
	// Recv: the stream window fills, the server's send blocks,
	// and the fan-out channel behind it has nowhere to drain.
	release := make(chan struct{})
	listenErrCh := make(chan error, 1)
	go func() {
		listenErrCh <- client.Listen(
			ctx, nil, 0,
			func(EntryMsg) error {
				<-release
				return nil
			},
		)
	}()

	waitForListeners(t, srv, 1)

	entries := []Entry{{
		ID:      "slow",
		Type:    "learning",
		Content: strings.Repeat("x", slowListenerPayload),
	}}
	var broadcasts int
	for srv.listeners.count() != 0 {
		if broadcasts == slowListenerBroadcastCap {
			t.Fatalf("listener still subscribed after %d "+
				"broadcasts to a client that never reads",
				broadcasts)
		}
		srv.listeners.broadcast(entries)
		broadcasts++
	}

	close(release)

	select {
	case listenErr := <-listenErrCh:
		if listenErr == nil {
			t.Fatal("Listen returned nil: the client cannot " +
				"tell a disconnect from a clean end of stream")
		}
		if code := status.Code(listenErr); code !=
			codes.ResourceExhausted {
			t.Errorf("Listen error code = %v (%v), want %v",
				code, listenErr, codes.ResourceExhausted)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Listen never returned after the server " +
			"disconnected its listener")
	}
}
