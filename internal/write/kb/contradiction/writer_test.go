//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package contradiction_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/contradiction"
)

func TestAppend_FirstRowGetsC001(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Contradictions)

	id, err := contradiction.Append(path, contradiction.Row{
		EVRefs:          []string{"EV-042", "EV-051"},
		Summary:         "widget scope disagreement",
		DemotionApplied: "EV-051 -> low",
		Status:          "resolved",
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "C-001" {
		t.Errorf("first id: want C-001; got %q", id)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if !strings.Contains(text, "| ID | Evidence |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "C-001") {
		t.Errorf("ID missing: %q", text)
	}
	if !strings.Contains(text, "EV-042, EV-051") {
		t.Errorf("evidence joined missing: %q", text)
	}
}

func TestAppend_AllocatesNextIDFromHighWater(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Contradictions)

	// Seed file with gaps: C-001, C-005, next must be C-006.
	seed := "| ID | Evidence | Summary | Demotion applied | Status |\n" +
		"|----|----------|---------|------------------|--------|\n" +
		"| C-001 | EV-001, EV-002 | s | d | resolved |\n" +
		"| C-005 | EV-003, EV-004 | s | d | open |\n"
	if err := os.WriteFile(path, []byte(seed), 0o600); err != nil {
		t.Fatal(err)
	}

	id, err := contradiction.Append(path, contradiction.Row{
		EVRefs: []string{"EV-009", "EV-010"}, Summary: "x",
		DemotionApplied: "y", Status: "open",
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "C-006" {
		t.Errorf("next id: want C-006; got %q", id)
	}
}

func TestAppend_SequentialIDsInSameProcess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Contradictions)

	var ids []string
	for i := 0; i < 3; i++ {
		id, err := contradiction.Append(path, contradiction.Row{
			EVRefs: []string{"EV-001", "EV-002"}, Summary: "x",
			DemotionApplied: "y", Status: "open",
		})
		if err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
		ids = append(ids, id)
	}
	want := []string{"C-001", "C-002", "C-003"}
	for i := range want {
		if ids[i] != want[i] {
			t.Errorf("id[%d]: want %s; got %s", i, want[i], ids[i])
		}
	}
}
