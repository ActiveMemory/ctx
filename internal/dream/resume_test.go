//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"reflect"
	"testing"
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// TestCrashResume simulates a pass that crashes mid-way: state was
// persisted only for the sources it completed. The next run must reload
// that committed state intact (SaveState writes atomically, so no torn
// file survives a crash) and the discipline clock must re-select exactly
// the work that was not finished — the sources never processed, plus any
// completed source whose content changed since — while skipping the ones
// already recorded unchanged. This is the spec's "next run resumes from
// the delta" / "no torn state" guarantee (specs/ctx-dream.md).
func TestCrashResume(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	dreamsDir := t.TempDir()

	hA := dream.Hash([]byte("alpha"))
	hB := dream.Hash([]byte("bravo"))
	hC := dream.Hash([]byte("charlie"))

	// The crashed pass completed A and B before dying on C; state is
	// persisted per completed item.
	partial := []dream.SourceState{
		{
			Path: "ideas/a.md", Hash: hA,
			Status: cfgDream.SourceActive, LastModified: time.Unix(0, 0).UTC(),
		},
		{
			Path: "ideas/b.md", Hash: hB,
			Status: cfgDream.SourceActive, LastModified: time.Unix(0, 0).UTC(),
		},
	}
	if err := dream.SaveState(dreamsDir, partial); err != nil {
		t.Fatalf("persist partial state: %v", err)
	}

	// Next run reloads the committed state — the atomic write means the
	// completed records survive the crash with no corruption.
	reloaded, err := dream.LoadState(dreamsDir)
	if err != nil {
		t.Fatalf("reload state: %v", err)
	}
	if len(reloaded) != 2 {
		t.Fatalf("expected 2 completed records to survive, got %d", len(reloaded))
	}

	// This pass sees all three ideas: A unchanged, B edited since the
	// crash, C never processed.
	current := map[string]string{
		"ideas/a.md": hA,                             // unchanged → skip
		"ideas/b.md": dream.Hash([]byte("bravo-v2")), // changed → re-triage
		"ideas/c.md": hC,                             // new → triage
	}

	got := dream.DeltaSelect(reloaded, current)
	want := []string{"ideas/b.md", "ideas/c.md"} // DeltaSelect sorts
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("resume delta = %v, want %v", got, want)
	}
}
