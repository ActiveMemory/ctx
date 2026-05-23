//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestFailoverClient_FirstPeerWorks verifies that the
// failover client connects to the first reachable peer.
func TestFailoverClient_FirstPeerWorks(t *testing.T) {
	_, _, adminTok := startTestServer(t)

	// Start a second server.
	dir := t.TempDir()
	store, storeErr := NewStore(dir)
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	srv := NewServer(store, adminTok)
	lis := listenRandom(t)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(func() { srv.GracefulStop() })

	addr := lis.Addr().String()

	// Register a client on the second server.
	regClient, dialErr := NewClient(addr, "")
	if dialErr != nil {
		t.Fatal(dialErr)
	}
	resp, regErr := regClient.Register(
		testCtx(), adminTok, "failover-proj",
	)
	if regErr != nil {
		t.Fatal(regErr)
	}
	_ = regClient.Close()

	// Failover client with the reachable peer first.
	client, foErr := newFailoverClient(
		[]string{addr}, resp.ClientToken,
	)
	if foErr != nil {
		t.Fatalf("newFailoverClient: %v", foErr)
	}
	defer func() { _ = client.Close() }()

	status, statusErr := client.Status(testCtx())
	if statusErr != nil {
		t.Fatalf("Status: %v", statusErr)
	}
	if status.TotalEntries != 0 {
		t.Errorf("want 0 entries, got %d",
			status.TotalEntries)
	}
}

// TestFailoverClient_SkipsBadPeer verifies that unreachable
// peers are skipped.
func TestFailoverClient_SkipsBadPeer(t *testing.T) {
	_, _, adminTok := startTestServer(t)

	dir := t.TempDir()
	store, _ := NewStore(dir)
	srv := NewServer(store, adminTok)
	lis := listenRandom(t)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(func() { srv.GracefulStop() })

	addr := lis.Addr().String()

	regClient, _ := NewClient(addr, "")
	resp, _ := regClient.Register(
		testCtx(), adminTok, "skip-proj",
	)
	_ = regClient.Close()

	// First peer is unreachable, second is good.
	client, foErr := newFailoverClient(
		[]string{"127.0.0.1:1", addr},
		resp.ClientToken,
	)
	if foErr != nil {
		t.Fatalf("expected fallback to work: %v", foErr)
	}
	_ = client.Close()
}

// TestFailoverClient_AllBad verifies error when no peer is
// reachable.
func TestFailoverClient_AllBad(t *testing.T) {
	_, foErr := newFailoverClient(
		[]string{"127.0.0.1:1", "127.0.0.1:2"},
		"bad-token",
	)
	if foErr == nil {
		t.Fatal("expected error when all peers bad")
	}
}

// TestFailoverClient_FailsFastOnAuthError verifies that an
// auth failure on the first reachable peer halts the
// failover walk: subsequent peers are not contacted, since
// the same token would fail there too.
//
// The reachable first peer is a real server that rejects
// the invalid bearer with codes.Unauthenticated. The
// second peer is `127.0.0.1:1` — an unrouted port that
// would surface a connection-class error (Unavailable) if
// the walk continued past the first auth failure. So an
// Unauthenticated return code proves the walk stopped at
// the first peer; an Unavailable return code would prove a
// regression (the walk cycled past auth and tried the
// second peer, where dial succeeded but the unrouted port
// produced a different error class).
func TestFailoverClient_FailsFastOnAuthError(t *testing.T) {
	_, conn, _ := startTestServer(t)
	addr := conn.Target()

	_, foErr := newFailoverClient(
		[]string{addr, "127.0.0.1:1"},
		"bogus-token-that-the-server-will-reject",
	)
	if foErr == nil {
		t.Fatal("expected auth error on first peer; got nil")
	}
	s, ok := status.FromError(foErr)
	if !ok {
		t.Fatalf("expected gRPC status error; got %T: %v", foErr, foErr)
	}
	if s.Code() != codes.Unauthenticated && s.Code() != codes.PermissionDenied {
		t.Errorf(
			"got code %s; want Unauthenticated or PermissionDenied "+
				"(if Unavailable, the walk cycled past auth into the "+
				"unrouted second peer — auth-fast-fail regression)",
			s.Code(),
		)
	}
}
