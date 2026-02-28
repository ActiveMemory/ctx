//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPauseMarkerPath(t *testing.T) {
	got := pauseMarkerPath("abc-123")
	want := filepath.Join(secureTempDir(), "ctx-paused-abc-123")
	if got != want {
		t.Errorf("pauseMarkerPath() = %q, want %q", got, want)
	}
}

func TestPaused_NoMarker(t *testing.T) {
	turns := paused("nonexistent-session-test-pause")
	if turns != 0 {
		t.Errorf("paused() with no marker = %d, want 0", turns)
	}
}

func TestPaused_IncrementsCounter(t *testing.T) {
	sessionID := "test-pause-increment"
	path := pauseMarkerPath(sessionID)
	writeCounter(path, 0) // create marker
	t.Cleanup(func() { _ = os.Remove(path) })

	// First call: 0 → 1
	turns := paused(sessionID)
	if turns != 1 {
		t.Errorf("paused() first call = %d, want 1", turns)
	}

	// Second call: 1 → 2
	turns = paused(sessionID)
	if turns != 2 {
		t.Errorf("paused() second call = %d, want 2", turns)
	}

	// Third call: 2 → 3
	turns = paused(sessionID)
	if turns != 3 {
		t.Errorf("paused() third call = %d, want 3", turns)
	}
}

func TestPaused_DoublePauseResets(t *testing.T) {
	sessionID := "test-pause-double"
	path := pauseMarkerPath(sessionID)
	writeCounter(path, 7) // simulate 7 turns
	t.Cleanup(func() { _ = os.Remove(path) })

	// Reset by writing 0
	writeCounter(path, 0)
	turns := paused(sessionID)
	if turns != 1 {
		t.Errorf("paused() after reset = %d, want 1", turns)
	}
}

func TestPausedMessage_NotPaused(t *testing.T) {
	got := pausedMessage(0)
	if got != "" {
		t.Errorf("pausedMessage(0) = %q, want empty", got)
	}
}

func TestPausedMessage_EarlyTurns(t *testing.T) {
	for _, n := range []int{1, 2, 3, 4, 5} {
		got := pausedMessage(n)
		if got != "ctx:paused" {
			t.Errorf("pausedMessage(%d) = %q, want %q", n, got, "ctx:paused")
		}
	}
}

func TestPausedMessage_LaterTurns(t *testing.T) {
	got := pausedMessage(6)
	want := "ctx:paused (6 turns) — resume with /ctx-resume"
	if got != want {
		t.Errorf("pausedMessage(6) = %q, want %q", got, want)
	}

	got = pausedMessage(100)
	want = "ctx:paused (100 turns) — resume with /ctx-resume"
	if got != want {
		t.Errorf("pausedMessage(100) = %q, want %q", got, want)
	}
}

func TestResumeWhenNotPaused(t *testing.T) {
	path := pauseMarkerPath("test-resume-noop")
	// Ensure no marker exists
	_ = os.Remove(path)

	// Remove on a non-existent file should not error
	removeErr := os.Remove(path)
	if removeErr != nil && !os.IsNotExist(removeErr) {
		t.Errorf("unexpected error removing non-existent marker: %v", removeErr)
	}
}

func TestResumeRemovesMarker(t *testing.T) {
	sessionID := "test-resume-removes"
	path := pauseMarkerPath(sessionID)
	writeCounter(path, 5)

	// Verify marker exists
	if _, statErr := os.Stat(path); statErr != nil {
		t.Fatalf("marker should exist: %v", statErr)
	}

	_ = os.Remove(path)

	// Verify marker is gone
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Error("marker should be removed after resume")
	}
}
