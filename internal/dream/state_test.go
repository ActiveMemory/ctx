//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	"github.com/ActiveMemory/ctx/internal/dream"
)

// writeFixture writes content to path, creating parent dirs. Shared by
// the dream engine tests.
func writeFixture(path, content string) error {
	if mkErr := os.MkdirAll(
		filepath.Dir(path), 0o755,
	); mkErr != nil {
		return mkErr
	}
	return os.WriteFile(path, []byte(content), 0o600)
}

// TestStateRoundTrip saves a state slice and loads it back unchanged.
func TestStateRoundTrip(t *testing.T) {
	dreamsDir := filepath.Join(t.TempDir(), "dreams")

	at := time.Date(2026, 6, 7, 2, 30, 0, 0, time.UTC)
	want := []dream.SourceState{
		{
			Path:         "ideas/a.md",
			Hash:         dream.Hash([]byte("alpha")),
			LastModified: at,
			Merit:        0.5,
			Status:       cfgDream.SourceActive,
		},
		{
			Path:         "ideas/b.md",
			Hash:         dream.Hash([]byte("beta")),
			LastModified: at,
			Merit:        0.9,
			Status:       cfgDream.SourcePromoted,
		},
	}

	if err := dream.SaveState(dreamsDir, want); err != nil {
		t.Fatalf("SaveState: %v", err)
	}
	got, err := dream.LoadState(dreamsDir)
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", got, want)
	}
}

// TestLoadStateMissing returns an empty slice when no state file exists.
func TestLoadStateMissing(t *testing.T) {
	dreamsDir := filepath.Join(t.TempDir(), "dreams")
	got, err := dream.LoadState(dreamsDir)
	if err != nil {
		t.Fatalf("LoadState: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("want empty, got %d entries", len(got))
	}
}

// TestDeltaSelect verifies the discipline clock selects new and changed
// sources and skips unchanged ones.
func TestDeltaSelect(t *testing.T) {
	prior := []dream.SourceState{
		{Path: "ideas/a.md", Hash: dream.Hash([]byte("alpha"))},
		{Path: "ideas/b.md", Hash: dream.Hash([]byte("beta"))},
	}
	current := map[string]string{
		"ideas/a.md": dream.Hash([]byte("alpha")),        // unchanged
		"ideas/b.md": dream.Hash([]byte("beta-changed")), // changed
		"ideas/c.md": dream.Hash([]byte("gamma")),        // new
	}

	got := dream.DeltaSelect(prior, current)
	want := []string{"ideas/b.md", "ideas/c.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DeltaSelect = %v, want %v", got, want)
	}
}

// TestDeltaSelectEmptyPriorAll selects every current source when there
// is no prior state (first-ever pass).
func TestDeltaSelectEmptyPriorAll(t *testing.T) {
	current := map[string]string{
		"ideas/a.md": dream.Hash([]byte("a")),
		"ideas/b.md": dream.Hash([]byte("b")),
	}
	got := dream.DeltaSelect(nil, current)
	want := []string{"ideas/a.md", "ideas/b.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DeltaSelect = %v, want %v", got, want)
	}
}
