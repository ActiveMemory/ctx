//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package glossary_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/glossary"
)

func TestAppend_CreatesHeaderOnFirstWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Glossary)

	row := glossary.Row{
		Term:         "widget",
		Definition:   "Smallest packaging boundary.",
		Confidence:   "high",
		EVRefs:       []string{"EV-042", "EV-043"},
		RelatedTerms: []string{"widget-contract"},
	}
	if err := glossary.Append(path, row); err != nil {
		t.Fatalf("Append: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if !strings.Contains(text, "| Term | Definition |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "widget") {
		t.Errorf("term missing: %q", text)
	}
	if !strings.Contains(text, "EV-042, EV-043") {
		t.Errorf("EV refs not joined: %q", text)
	}
}

func TestAppend_SecondCallDoesNotRewriteHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Glossary)

	for i, term := range []string{"alpha", "beta"} {
		row := glossary.Row{
			Term: term, Definition: "def", Confidence: "low",
			EVRefs: []string{"EV-001"},
		}
		if err := glossary.Append(path, row); err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if c := strings.Count(text, "| Term | Definition |"); c != 1 {
		t.Errorf("header written %d times; want 1", c)
	}
	if !strings.Contains(text, "alpha") || !strings.Contains(text, "beta") {
		t.Errorf("both rows missing: %q", text)
	}
}

func TestAppend_EscapesPipesAndNewlines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.Glossary)

	row := glossary.Row{
		Term:       "tricky",
		Definition: "has a | pipe\nand newline",
		Confidence: "medium",
		EVRefs:     []string{"EV-099"},
	}
	if err := glossary.Append(path, row); err != nil {
		t.Fatalf("Append: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	text := string(got)
	if strings.Contains(text, "has a | pipe") {
		t.Errorf("unescaped pipe in cell: %q", text)
	}
	if !strings.Contains(text, `\|`) {
		t.Errorf("expected escaped pipe: %q", text)
	}
	// The row should be one line (excluding header lines).
	rows := strings.Split(strings.TrimRight(text, "\n"), "\n")
	// header (2 lines) + 1 row = 3 lines
	if len(rows) != 3 {
		t.Errorf("expected 3 lines; got %d: %q", len(rows), text)
	}
}
