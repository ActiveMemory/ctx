//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package rc

import (
	"testing"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
)

// TestDreamDefaults verifies the dream accessors fall back to the
// config/dream defaults when the dream section is absent.
func TestDreamDefaults(t *testing.T) {
	declareContext(t, "")

	if DreamEnabled() {
		t.Error("DreamEnabled() = true, want false (opt-in default)")
	}
	if got := DreamMode(); got != cfgDream.ModeDiscipline {
		t.Errorf("DreamMode() = %q, want %q", got, cfgDream.ModeDiscipline)
	}
	if got := DreamMax(); got != cfgDream.DefaultMax {
		t.Errorf("DreamMax() = %d, want %d", got, cfgDream.DefaultMax)
	}
	if got := DreamBudget(); got != cfgDream.DefaultBudget {
		t.Errorf("DreamBudget() = %d, want %d", got, cfgDream.DefaultBudget)
	}
	if got := DreamQuietMinutes(); got != cfgDream.DefaultQuietMinutes {
		t.Errorf(
			"DreamQuietMinutes() = %d, want %d",
			got, cfgDream.DefaultQuietMinutes,
		)
	}
	if got := DreamCadence(); got != "" {
		t.Errorf("DreamCadence() = %q, want empty", got)
	}
	if got := DreamModel(); got != "" {
		t.Errorf("DreamModel() = %q, want empty", got)
	}
	if got := DreamExecutor(); got != "" {
		t.Errorf("DreamExecutor() = %q, want empty", got)
	}
}

// TestDreamConfigured verifies the dream accessors read explicit
// .ctxrc values, including the creative mode.
func TestDreamConfigured(t *testing.T) {
	declareContext(t, `dream:
  enabled: true
  mode: creative
  max: 12
  cadence: "30 2 * * *"
  quiet_minutes: 90
  model: opus
  budget: 7
  executor: "my-runner --headless"
`)

	if !DreamEnabled() {
		t.Error("DreamEnabled() = false, want true")
	}
	if got := DreamMode(); got != cfgDream.ModeCreative {
		t.Errorf("DreamMode() = %q, want %q", got, cfgDream.ModeCreative)
	}
	if got := DreamMax(); got != 12 {
		t.Errorf("DreamMax() = %d, want 12", got)
	}
	if got := DreamCadence(); got != "30 2 * * *" {
		t.Errorf("DreamCadence() = %q, want cron string", got)
	}
	if got := DreamQuietMinutes(); got != 90 {
		t.Errorf("DreamQuietMinutes() = %d, want 90", got)
	}
	if got := DreamModel(); got != "opus" {
		t.Errorf("DreamModel() = %q, want opus", got)
	}
	if got := DreamBudget(); got != 7 {
		t.Errorf("DreamBudget() = %d, want 7", got)
	}
	if got := DreamExecutor(); got != "my-runner --headless" {
		t.Errorf("DreamExecutor() = %q, want runner string", got)
	}
}
