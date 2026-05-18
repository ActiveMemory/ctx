//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package relationship_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/relationship"
)

func TestAppend_CreatesHeaderOnFirstWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.RelationshipMap)

	row := relationship.Row{
		From: "widgets", To: "EV-042",
		Kind:    "depends-on",
		Summary: "Widget topic page leans on EV-042 for scope.",
	}
	if err := relationship.Append(path, row); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if !strings.Contains(text, "| From | To | Kind | Summary |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "widgets") || !strings.Contains(text, "EV-042") {
		t.Errorf("endpoints missing: %q", text)
	}
	if !strings.Contains(text, "depends-on") {
		t.Errorf("kind missing: %q", text)
	}
}

func TestAppend_TwoRowsSingleHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.RelationshipMap)

	for _, k := range []string{"depends-on", "refines"} {
		if err := relationship.Append(path, relationship.Row{
			From: "a", To: "b", Kind: k, Summary: "x",
		}); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	got, _ := os.ReadFile(path)
	if c := strings.Count(string(got), "| From | To | Kind | Summary |"); c != 1 {
		t.Errorf("header count: want 1; got %d", c)
	}
}
