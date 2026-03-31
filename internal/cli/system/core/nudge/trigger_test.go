//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

import (
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/event"
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
