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

		// Write two stats entries — the second should win.
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
