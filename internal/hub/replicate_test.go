//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"bytes"
	"net"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// captureWarnings redirects the warn sink to a buffer for the
// duration of the test.
func captureWarnings(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	restore := logWarn.SetSink(&buf)
	t.Cleanup(restore)
	return &buf
}

// startMaster serves a hub over a fresh store on a random port
// and returns the store, the address, and the admin token.
func startMaster(t *testing.T) (*Store, string, string) {
	t.Helper()
	store, storeErr := NewStore(t.TempDir())
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	adminTok, tokErr := GenerateAdminToken()
	if tokErr != nil {
		t.Fatal(tokErr)
	}
	srv := NewServer(store, adminTok)
	lis := listenRandom(t)
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.GracefulStop)
	return store, lis.Addr().String(), adminTok
}

// registerClient registers a project and returns its token.
func registerClient(t *testing.T, addr, adminTok string) string {
	t.Helper()
	client, dialErr := NewClient(addr, "")
	if dialErr != nil {
		t.Fatal(dialErr)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			t.Log(cerr)
		}
	}()
	reg, regErr := client.Register(
		testCtx(), adminTok, "replicate-test",
	)
	if regErr != nil {
		t.Fatal(regErr)
	}
	return reg.ClientToken
}

// seedEntries appends n entries to the master store.
func seedEntries(t *testing.T, store *Store, n int) {
	t.Helper()
	entries := make([]Entry, n)
	for i := range entries {
		entries[i] = Entry{
			ID:        "e" + string(rune('a'+i)),
			Type:      "decision",
			Content:   "replicated content",
			Origin:    "replicate-test",
			Timestamp: time.Now(),
		}
	}
	if _, appendErr := store.Append(entries); appendErr != nil {
		t.Fatal(appendErr)
	}
}

func TestReplicateOnce_WarnsOnDialError(t *testing.T) {
	buf := captureWarnings(t)
	follower, storeErr := NewStore(t.TempDir())
	if storeErr != nil {
		t.Fatal(storeErr)
	}

	// grpc.NewClient is lazy for almost every bad target, but
	// a control character fails URL parsing at construction —
	// the one eager failure mode, and exactly what a corrupted
	// peer config would produce.
	replicateOnce(testCtx(), "\x00", follower, "tok")

	if !strings.Contains(buf.String(), "hub replicate dial") {
		t.Errorf("missing dial warning, got: %q", buf.String())
	}
}

func TestReplicateOnce_WarnsOnTransportError(t *testing.T) {
	buf := captureWarnings(t)
	follower, storeErr := NewStore(t.TempDir())
	if storeErr != nil {
		t.Fatal(storeErr)
	}

	// A well-formed address nobody listens on. Client and
	// stream construction are lazy, so the connection failure
	// surfaces at whichever stage first touches the wire —
	// any of the replication warnings is correct.
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}
	addr := lis.Addr().String()
	if closeErr := lis.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}

	replicateOnce(testCtx(), addr, follower, "tok")

	if !strings.Contains(buf.String(), "hub replicate") {
		t.Errorf(
			"missing transport warning, got: %q", buf.String(),
		)
	}
}

func TestReplicateOnce_CleanReplicationDoesNotWarn(t *testing.T) {
	master, addr, adminTok := startMaster(t)
	token := registerClient(t, addr, adminTok)
	seedEntries(t, master, 2)

	follower, storeErr := NewStore(t.TempDir())
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	buf := captureWarnings(t)

	replicateOnce(testCtx(), addr, follower, token)

	if buf.Len() != 0 {
		t.Errorf(
			"clean replication must not warn, got: %q",
			buf.String(),
		)
	}
	if _, lastSeq := follower.lastSequence(); lastSeq != 2 {
		t.Errorf("follower lastSeq = %d, want 2", lastSeq)
	}
}

func TestReplicateOnce_KeepsConsumingAfterAppendError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission semantics differ on windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission checks")
	}
	master, addr, adminTok := startMaster(t)
	token := registerClient(t, addr, adminTok)
	seedEntries(t, master, 2)

	followerDir := t.TempDir()
	follower, storeErr := NewStore(followerDir)
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	// Make every append fail: the store opens its files per
	// call, so an access-denied directory rejects each write.
	if chmodErr := os.Chmod(followerDir, 0); chmodErr != nil {
		t.Fatal(chmodErr)
	}
	t.Cleanup(func() {
		if chmodErr := os.Chmod(
			followerDir, fs.PermExec,
		); chmodErr != nil {
			t.Log(chmodErr)
		}
	})
	buf := captureWarnings(t)

	replicateOnce(testCtx(), addr, follower, token)

	// Both entries must have been attempted: an append failure
	// is warned per entry and must not abort the stream.
	got := strings.Count(buf.String(), "hub replicate append")
	if got != 2 {
		t.Errorf(
			"append warnings = %d, want 2 (loop must keep "+
				"consuming); output: %q",
			got, buf.String(),
		)
	}
}
