//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

func TestFanOut_SubscribeAndBroadcast(t *testing.T) {
	fo := newFanOut()

	ch1 := fo.subscribe()
	ch2 := fo.subscribe()

	if fo.count() != 2 {
		t.Fatalf("want 2 subs, got %d", fo.count())
	}

	entries := []Entry{
		{ID: "x", Content: "test"},
	}
	fo.broadcast(entries)

	select {
	case got := <-ch1:
		if got[0].ID != "x" {
			t.Errorf("ch1: want ID 'x', got %q", got[0].ID)
		}
	case <-time.After(time.Second):
		t.Fatal("ch1: timeout")
	}

	select {
	case got := <-ch2:
		if got[0].ID != "x" {
			t.Errorf("ch2: want ID 'x', got %q", got[0].ID)
		}
	case <-time.After(time.Second):
		t.Fatal("ch2: timeout")
	}
}

func TestFanOut_Unsubscribe(t *testing.T) {
	fo := newFanOut()
	ch := fo.subscribe()
	fo.unsubscribe(ch)

	if fo.count() != 0 {
		t.Errorf("want 0 subs after unsubscribe, got %d",
			fo.count())
	}
}

func TestFanOut_BroadcastToNone(t *testing.T) {
	fo := newFanOut()
	// Should not panic.
	fo.broadcast([]Entry{{ID: "noop"}})
}

func TestFanOut_DisconnectsSlowListener(t *testing.T) {
	// The disconnect warns on stderr; keep test output clean.
	restore := logWarn.SetSink(io.Discard)
	defer restore()

	fo := newFanOut()
	slow := fo.subscribe()

	// Never read from slow. One broadcast past the buffer has
	// nowhere to go, so the listener is disconnected.
	for i := 0; i < fanOutBuffer+1; i++ {
		fo.broadcast([]Entry{{ID: fmt.Sprintf("e%d", i)}})
	}

	if got := fo.count(); got != 0 {
		t.Errorf("count = %d, want 0 after disconnect", got)
	}
	if got := fo.droppedCount(); got != 1 {
		t.Errorf("droppedCount = %d, want 1", got)
	}

	// Drain the buffered entries, then observe the close.
	deadline := time.After(time.Second)
	for i := 0; i <= fanOutBuffer; i++ {
		select {
		case _, ok := <-slow:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("disconnected channel never closed")
		}
	}
	select {
	case _, ok := <-slow:
		if ok {
			t.Fatal("channel still open after disconnect")
		}
	case <-deadline:
		t.Fatal("disconnected channel never closed")
	}
}

// TestFanOut_UnsubscribeAfterDisconnect covers the collision
// between the two paths that close a listener channel: broadcast
// disconnecting a slow listener, and the Listen stream's deferred
// unsubscribe on the way out. The second close must be a no-op —
// closing an already-closed channel panics, and the gRPC server
// has no recovery interceptor to catch it.
func TestFanOut_UnsubscribeAfterDisconnect(t *testing.T) {
	restore := logWarn.SetSink(io.Discard)
	defer restore()

	fo := newFanOut()
	slow := fo.subscribe()

	for i := 0; i <= fanOutBuffer; i++ {
		fo.broadcast([]Entry{{ID: fmt.Sprintf("e%d", i)}})
	}
	if got := fo.count(); got != 0 {
		t.Fatalf("count = %d, want 0 after disconnect", got)
	}

	// Would panic without the membership guard.
	fo.unsubscribe(slow)
	// Idempotent for any number of callers.
	fo.unsubscribe(slow)

	if got := fo.droppedCount(); got != 1 {
		t.Errorf("droppedCount = %d, want 1 (unsubscribe must "+
			"not count as a drop)", got)
	}
}

func TestFanOut_DroppedCountStartsAtZero(t *testing.T) {
	fo := newFanOut()
	ch := fo.subscribe()
	fo.broadcast([]Entry{{ID: "x"}})
	<-ch

	if got := fo.droppedCount(); got != 0 {
		t.Errorf("droppedCount = %d, want 0 for a healthy listener",
			got)
	}
}

// TestFanOut_DroppedCountRaceWithBroadcast exercises the read
// path the Status RPC handler uses: droppedCount from another
// goroutine while broadcast is disconnecting listeners. Run
// under -race, it fails if the counter stops being atomic.
func TestFanOut_DroppedCountRaceWithBroadcast(t *testing.T) {
	restore := logWarn.SetSink(io.Discard)
	defer restore()

	fo := newFanOut()

	var wg sync.WaitGroup
	done := make(chan struct{})

	// Readers stand in for concurrent Status RPC handlers.
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
					_ = fo.droppedCount()
				}
			}
		}()
	}

	// Each round subscribes a listener that never drains, then
	// overflows it so broadcast disconnects it.
	for round := 0; round < 20; round++ {
		fo.subscribe()
		for i := 0; i <= fanOutBuffer; i++ {
			fo.broadcast([]Entry{{ID: fmt.Sprintf("r%d-%d", round, i)}})
		}
	}

	close(done)
	wg.Wait()

	if got := fo.droppedCount(); got != 20 {
		t.Errorf("droppedCount = %d, want 20", got)
	}
}
