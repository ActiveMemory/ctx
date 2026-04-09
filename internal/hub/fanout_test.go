//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hub

import (
	"testing"
	"time"
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
