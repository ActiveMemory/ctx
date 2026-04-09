//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"testing"
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
	client, foErr := NewFailoverClient(
		[]string{addr}, resp.ClientToken,
	)
	if foErr != nil {
		t.Fatalf("NewFailoverClient: %v", foErr)
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
	client, foErr := NewFailoverClient(
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
	_, foErr := NewFailoverClient(
		[]string{"127.0.0.1:1", "127.0.0.1:2"},
		"bad-token",
	)
	if foErr == nil {
		t.Fatal("expected error when all peers bad")
	}
}
