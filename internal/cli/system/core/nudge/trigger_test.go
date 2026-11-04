//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

func TestEvaluateTrigger(t *testing.T) {
	tests := []struct {
		name            string
		pct             int
		checkpointFired bool
		wantEvent       string
		wantCheckpoint  bool
		wantWindow      bool
	}{
		{
			name:      "pct 0 no nudge",
			pct:       0,
			wantEvent: event.Silent,
		},
		{
			name:      "pct 30 no nudge",
			pct:       30,
			wantEvent: event.Silent,
		},
		{
			name:      "pct 59 no nudge",
			pct:       59,
			wantEvent: event.Silent,
		},
		{
			name:           "pct 60 fires checkpoint",
			pct:            60,
			wantEvent:      event.Checkpoint,
			wantCheckpoint: true,
		},
		{
			name:           "pct 75 fires checkpoint",
			pct:            75,
			wantEvent:      event.Checkpoint,
			wantCheckpoint: true,
		},
		{
			name:            "pct 65 checkpoint already fired is silent",
			pct:             65,
			checkpointFired: true,
			wantEvent:       event.Silent,
		},
		{
			name:            "pct 89 checkpoint already fired is silent",
			pct:             89,
			checkpointFired: true,
			wantEvent:       event.Silent,
		},
		{
			name:       "pct 90 fires window warning",
			pct:        90,
			wantEvent:  event.WindowWarning,
			wantWindow: true,
		},
		{
			name:       "pct 95 fires window warning",
			pct:        95,
			wantEvent:  event.WindowWarning,
			wantWindow: true,
		},
		{
			name:       "pct 100 fires window warning",
			pct:        100,
			wantEvent:  event.WindowWarning,
			wantWindow: true,
		},
		{
			name:            "pct 90 with checkpoint fired still fires window",
			pct:             90,
			checkpointFired: true,
			wantEvent:       event.WindowWarning,
			wantWindow:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateTrigger(tt.pct, tt.checkpointFired)
			if got.Event != tt.wantEvent {
				t.Errorf("Event = %q, want %q", got.Event, tt.wantEvent)
			}
			if got.Checkpoint != tt.wantCheckpoint {
				t.Errorf("Checkpoint = %v, want %v",
					got.Checkpoint, tt.wantCheckpoint)
			}
			if got.Window != tt.wantWindow {
				t.Errorf("Window = %v, want %v", got.Window, tt.wantWindow)
			}
		})
	}
}

func TestEmitEoCCWarning_OneShotAndPrecedence(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(origDir)
		rc.Reset()
	})

	ctxDir := filepath.Join(tmpDir, ".context")
	for _, sub := range []string{"state", "logs"} {
		if err := os.MkdirAll(filepath.Join(ctxDir, sub), 0o750); err != nil {
			t.Fatal(err)
		}
	}
	state.SetDirForTest(filepath.Join(ctxDir, "state"))
	t.Cleanup(func() { state.SetDirForTest("") })
	testctx.Declare(t, tmpDir)
	rc.Reset()

	logFile := filepath.Join(ctxDir, "logs", "check-context-size.log")
	box, emitErr := EmitEoCCWarning(logFile, "sess-eocc", 5, 32000, 30000)
	if emitErr != nil {
		t.Fatalf("EmitEoCCWarning() first call error = %v", emitErr)
	}
	if !strings.Contains(box, "EoCC Compaction Threshold") {
		t.Errorf("expected EoCC box title in %q", box)
	}
	if !strings.Contains(box, "ctx handover write") {
		t.Errorf("expected handover instruction in %q", box)
	}

	// Second call in the same session should be suppressed by the one-shot guard.
	secondBox, secondErr := EmitEoCCWarning(logFile, "sess-eocc", 6, 35000, 30000)
	if secondErr != nil {
		t.Fatalf("EmitEoCCWarning() second call error = %v", secondErr)
	}
	if secondBox != "" {
		t.Errorf("expected empty box on second call, got %q", secondBox)
	}
}

