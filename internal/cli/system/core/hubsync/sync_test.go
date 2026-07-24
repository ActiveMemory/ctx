//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hubsync_test

import (
	"bytes"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/lookup"
	connectCfg "github.com/ActiveMemory/ctx/internal/cli/connection/core/config"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/hubsync"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/hub"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

// TestMain initializes the embedded text-asset lookup so error
// strings rendered into warnings resolve their DescKey text.
func TestMain(m *testing.M) {
	lookup.Init()
	os.Exit(m.Run())
}

// declareContext positions the test in a temp project with a
// materialized .context/ directory and returns its path.
func declareContext(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	ctxDir := testctx.Declare(t, dir)
	if mkErr := os.MkdirAll(ctxDir, fs.PermExec); mkErr != nil {
		t.Fatal(mkErr)
	}
	return ctxDir
}

// captureWarnings redirects the warn sink to a buffer for the
// duration of the test.
func captureWarnings(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	restore := logWarn.SetSink(&buf)
	t.Cleanup(restore)
	return &buf
}

// saveConnectConfig generates a global encryption key under the
// test HOME and persists a connection config pointing at addr.
func saveConnectConfig(t *testing.T, addr, token string) {
	t.Helper()
	key, keyErr := crypto.GenerateKey()
	if keyErr != nil {
		t.Fatal(keyErr)
	}
	keyPath := crypto.GlobalKeyPath()
	if mkErr := os.MkdirAll(
		filepath.Dir(keyPath), fs.PermKeyDir,
	); mkErr != nil {
		t.Fatal(mkErr)
	}
	if saveErr := crypto.SaveKey(keyPath, key); saveErr != nil {
		t.Fatal(saveErr)
	}
	if cfgErr := connectCfg.Save(connectCfg.Config{
		HubAddr: addr,
		Token:   token,
	}); cfgErr != nil {
		t.Fatal(cfgErr)
	}
}

// startHub serves a hub with the given store on a random port
// and returns its address and a registered client token.
func startHub(t *testing.T, store *hub.Store) (string, string) {
	t.Helper()
	adminTok, tokErr := hub.GenerateAdminToken()
	if tokErr != nil {
		t.Fatal(tokErr)
	}
	srv := hub.NewServer(store, adminTok)
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.GracefulStop)

	addr := lis.Addr().String()
	client, dialErr := hub.NewClient(addr, "")
	if dialErr != nil {
		t.Fatal(dialErr)
	}
	defer func() {
		if cerr := client.Close(); cerr != nil {
			t.Log(cerr)
		}
	}()
	reg, regErr := client.Register(
		t.Context(), adminTok, "hubsync-test",
	)
	if regErr != nil {
		t.Fatal(regErr)
	}
	return addr, reg.ClientToken
}

func TestSync_WarnsOnLoadError(t *testing.T) {
	declareContext(t)
	buf := captureWarnings(t)

	if got := hubsync.Sync(""); got != "" {
		t.Errorf("Sync = %q, want empty on load failure", got)
	}
	if !strings.Contains(
		buf.String(), "hubsync: load connection config:",
	) {
		t.Errorf("missing load warning, got: %q", buf.String())
	}
}

func TestSync_WarnsOnDialError(t *testing.T) {
	declareContext(t)
	// grpc.NewClient is lazy for almost every bad target, but
	// a control character fails URL parsing at construction —
	// the one eager failure mode, and exactly what a corrupted
	// connect config would produce.
	saveConnectConfig(t, "\x00", "tok")
	buf := captureWarnings(t)

	if got := hubsync.Sync(""); got != "" {
		t.Errorf("Sync = %q, want empty on dial failure", got)
	}
	if !strings.Contains(buf.String(), "hubsync: dial") {
		t.Errorf("missing dial warning, got: %q", buf.String())
	}
}

func TestSync_WarnsOnPullError(t *testing.T) {
	declareContext(t)
	// A well-formed address nobody listens on: client
	// construction is lazy, so the failure surfaces at the
	// Sync RPC.
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}
	addr := lis.Addr().String()
	if closeErr := lis.Close(); closeErr != nil {
		t.Fatal(closeErr)
	}
	saveConnectConfig(t, addr, "tok")
	buf := captureWarnings(t)

	if got := hubsync.Sync(""); got != "" {
		t.Errorf("Sync = %q, want empty on pull failure", got)
	}
	if !strings.Contains(buf.String(), "hubsync: sync from") {
		t.Errorf("missing pull warning, got: %q", buf.String())
	}
}

func TestSync_NoWarnOnEmptyResult(t *testing.T) {
	ctxDir := declareContext(t)
	store, storeErr := hub.NewStore(
		filepath.Join(ctxDir, "hub-data"),
	)
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	addr, token := startHub(t, store)
	saveConnectConfig(t, addr, token)
	buf := captureWarnings(t)

	if got := hubsync.Sync(""); got != "" {
		t.Errorf("Sync = %q, want empty for zero entries", got)
	}
	if buf.Len() != 0 {
		t.Errorf(
			"empty result must not warn, got: %q", buf.String(),
		)
	}
}

// blackHoleListener accepts TCP connections but never speaks
// gRPC, so a dial succeeds while every RPC blocks until its
// deadline. It models a hung-but-reachable hub. Accepted
// connections and the listener are closed at test cleanup.
func blackHoleListener(t *testing.T) string {
	t.Helper()
	lis, lisErr := net.Listen("tcp", "127.0.0.1:0")
	if lisErr != nil {
		t.Fatal(lisErr)
	}
	var mu sync.Mutex
	var conns []net.Conn
	t.Cleanup(func() {
		_ = lis.Close()
		mu.Lock()
		for _, c := range conns {
			_ = c.Close()
		}
		mu.Unlock()
	})
	go func() {
		for {
			conn, acceptErr := lis.Accept()
			if acceptErr != nil {
				return
			}
			mu.Lock()
			conns = append(conns, conn)
			mu.Unlock()
		}
	}()
	return lis.Addr().String()
}

func TestSync_WarnsOnHungHub(t *testing.T) {
	declareContext(t)
	// Capture warnings before standing up the black-hole listener.
	// Cleanup is LIFO, so registering the warn-sink restore first
	// makes it run last — after the listener's connections close. On
	// the failure path (deadline not enforced), that ordering lets a
	// leaked Sync goroutine's late warning land in this test's buffer
	// rather than the restored process-global sink of the next test.
	buf := captureWarnings(t)
	// A hub that accepts the connection but never responds. The
	// RPC has no inherent deadline, so without HubSyncTimeout this
	// would hang the session-start hook indefinitely. Shrink the
	// bound so the test exercises the deadline in milliseconds.
	saveConnectConfig(t, blackHoleListener(t), "tok")
	restore := hubsync.SetSyncTimeoutForTest(200 * time.Millisecond)
	t.Cleanup(restore)

	done := make(chan string, 1)
	go func() { done <- hubsync.Sync("") }()

	select {
	case got := <-done:
		if got != "" {
			t.Errorf("Sync = %q, want empty on deadline", got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Sync did not return: the deadline was not enforced")
	}
	if !strings.Contains(buf.String(), "hubsync: sync from") {
		t.Errorf("missing pull warning, got: %q", buf.String())
	}
}
