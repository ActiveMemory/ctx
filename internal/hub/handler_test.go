//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// listenTestServer returns a Server over an empty store. Serve is
// never called: the Listen handler is driven directly so the test
// owns the send callback and can stall it.
func listenTestServer(t *testing.T) *Server {
	t.Helper()

	store, storeErr := NewStore(t.TempDir())
	if storeErr != nil {
		t.Fatal(storeErr)
	}
	adminTok, tokErr := GenerateAdminToken()
	if tokErr != nil {
		t.Fatal(tokErr)
	}
	return NewServer(store, adminTok)
}

// waitForListeners blocks until the server's subscriber count
// reaches want, failing the test if it does not settle in time.
func waitForListeners(t *testing.T, srv *Server, want uint32) {
	t.Helper()

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if srv.listeners.count() == want {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("listener count = %d, want %d",
		srv.listeners.count(), want)
}

// TestListenEntries_SlowListenerEndsStream pins the whole
// disconnect path end to end. A stalled send lets the fan-out
// buffer fill; broadcast then disconnects the listener and closes
// its channel. The handler must notice the close and return
// errSlowListener rather than spin on a channel that is always
// receivable, and its deferred unsubscribe must survive the
// channel broadcast already closed — a second close panics, and
// the gRPC server carries no recovery interceptor, so the panic
// would take the hub daemon down.
func TestListenEntries_SlowListenerEndsStream(t *testing.T) {
	restore := logWarn.SetSink(io.Discard)
	defer restore()

	srv := listenTestServer(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	release := make(chan struct{})
	var delivered atomic.Int64
	done := make(chan error, 1)

	go func() {
		done <- srv.listenEntries(
			&ListenRequest{},
			func(*EntryMsg) error {
				<-release
				delivered.Add(1)
				return nil
			},
			ctx,
		)
	}()

	waitForListeners(t, srv, 1)

	// One broadcast past the buffer, plus one for the entry the
	// handler takes off the channel before stalling in send.
	for i := 0; i <= fanOutBuffer+1; i++ {
		srv.listeners.broadcast(
			[]Entry{{ID: fmt.Sprintf("e%d", i)}},
		)
	}

	waitForListeners(t, srv, 0)
	if got := srv.listeners.droppedCount(); got != 1 {
		t.Fatalf("droppedCount = %d, want 1", got)
	}

	close(release)

	select {
	case listenErr := <-done:
		if !errors.Is(listenErr, errSlowListener) {
			t.Fatalf("listenEntries = %v, want errSlowListener",
				listenErr)
		}
		if code := status.Code(listenErr); code !=
			codes.ResourceExhausted {
			t.Errorf("status code = %v, want %v",
				code, codes.ResourceExhausted)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("listenEntries never returned after its " +
			"listener was disconnected")
	}

	// Whatever the buffer still held at disconnect time is sent
	// before the stream ends: the disconnect costs the client the
	// stream, not the entries already handed to it.
	if got := delivered.Load(); got < fanOutBuffer {
		t.Errorf("delivered %d entries before the disconnect, "+
			"want at least %d", got, fanOutBuffer)
	}
}

// TestListenEntries_ContextCancelEndsStream keeps the healthy
// shutdown path honest: a cancelled stream returns nil, and the
// deferred unsubscribe closes a channel that is still live.
func TestListenEntries_ContextCancelEndsStream(t *testing.T) {
	srv := listenTestServer(t)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)

	go func() {
		done <- srv.listenEntries(
			&ListenRequest{},
			func(*EntryMsg) error { return nil },
			ctx,
		)
	}()

	waitForListeners(t, srv, 1)
	cancel()

	select {
	case listenErr := <-done:
		if listenErr != nil {
			t.Fatalf("listenEntries = %v, want nil on cancel",
				listenErr)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("listenEntries never returned after cancel")
	}

	waitForListeners(t, srv, 0)
	if got := srv.listeners.droppedCount(); got != 0 {
		t.Errorf("droppedCount = %d, want 0 for a clean exit", got)
	}
}
