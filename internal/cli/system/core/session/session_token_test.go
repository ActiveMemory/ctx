//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/claude"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/stats"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/testutil/testctx"
)

func TestModelContextWindow(t *testing.T) {
	tests := []struct {
		name  string
		model string
		want  int
	}{
		{name: "empty model", model: "", want: 0},
		{name: "unknown model", model: "gpt-4", want: 0},
		{
			name:  "claude opus is always 1M",
			model: "claude-opus-4-6-20260205",
			want:  claude.ContextWindow1M,
		},
		{
			name:  "claude sonnet is 200k",
			model: "claude-sonnet-4-6-20260205",
			want:  rc.DefaultContextWindow,
		},
		{
			name:  "claude with 1m suffix",
			model: "claude-opus-4-6[1m]",
			want:  claude.ContextWindow1M,
		},
		{
			name:  "claude with 1M uppercase",
			model: "claude-sonnet-4-6[1M]",
			want:  claude.ContextWindow1M,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ModelContextWindow(tt.model)
			if got != tt.want {
				t.Errorf("ModelContextWindow(%q) = %d, want %d", tt.model, got, tt.want)
			}
		})
	}
}

func TestLatestSessionPct(t *testing.T) {
	// Set up a temp state dir.
	tmpDir := t.TempDir()
	state.SetDirForTest(tmpDir)
	t.Cleanup(func() { state.SetDirForTest("") })

	sessionID := "test-session-pct"

	t.Run("no stats file returns 0", func(t *testing.T) {
		got := LatestPct("nonexistent-session")
		if got != 0 {
			t.Errorf("LatestPct(nonexistent) = %d, want 0", got)
		}
	})

	t.Run("reads latest pct from stats JSONL", func(t *testing.T) {
		path := filepath.Join(tmpDir, stats.FilePrefix+sessionID+file.ExtJSONL)

		// Write two stats entries; the second should win.
		entries := []entity.Stats{
			{Pct: 5, Prompt: 1, Event: "silent"},
			{Pct: 12, Prompt: 2, Event: "silent"},
		}
		f, createErr := os.Create(path)
		if createErr != nil {
			t.Fatal(createErr)
		}
		for _, e := range entries {
			data, marshalErr := json.Marshal(e)
			if marshalErr != nil {
				t.Fatal(marshalErr)
			}
			if _, writeErr := f.Write(append(data, '\n')); writeErr != nil {
				t.Fatal(writeErr)
			}
		}
		if closeErr := f.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}

		got := LatestPct(sessionID)
		if got != 12 {
			t.Errorf("LatestPct(%q) = %d, want 12", sessionID, got)
		}
	})
}

// TestFindJSONLPathDoesNotMaterializeContext verifies that calling
// FindJSONLPath in a project that has not run "ctx init" does not
// silently create a phantom .context/ (or .context/state/) directory
// as a side effect of cache writeback.
//
// Provenance.Emit is intentionally unconditional, so it must be safe
// to call from any hook regardless of init state.
func TestFindJSONLPathDoesNotMaterializeContext(t *testing.T) {
	tmpDir := t.TempDir()
	t.Chdir(tmpDir)
	ctxPath := testctx.Declare(t, tmpDir)

	// Sanity: the .context/ dir does not exist yet.
	if _, statErr := os.Stat(ctxPath); !os.IsNotExist(statErr) {
		t.Fatalf("precondition: %s should not exist; statErr=%v",
			ctxPath, statErr)
	}

	path, err := FindJSONLPath("any-session-id")
	if err != nil {
		t.Fatalf("FindJSONLPath returned error: %v", err)
	}
	if path != "" {
		t.Errorf("FindJSONLPath returned %q, want empty (uninitialized project)",
			path)
	}

	// The .context/ directory must NOT have been materialized.
	if _, statErr := os.Stat(ctxPath); !os.IsNotExist(statErr) {
		t.Errorf("FindJSONLPath materialized %s in an uninitialized project; statErr=%v",
			ctxPath, statErr)
	}
}
