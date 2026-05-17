//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/decision"
)

func TestAppend_FirstRowGetsDD001(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.DomainDecisions)

	id, err := decision.Append(path, decision.Row{
		Date:         "2026-05-16",
		Context:      "Two sources disagreed on bundle opacity.",
		Decision:     "Treat widget bundles as opaque to the consumer.",
		Rationale:    "Opacity preserves producer's freedom to evolve internals.",
		Consequence:  "widget-composition rewritten; glossary widget gains opacity xref.",
		SupportingEV: []string{"EV-042", "EV-043"},
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "DD-001" {
		t.Errorf("first id: want DD-001; got %q", id)
	}

	got, _ := os.ReadFile(path)
	text := string(got)
	if !strings.Contains(text, "| ID | Date |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "DD-001") {
		t.Errorf("ID missing: %q", text)
	}
	if !strings.Contains(text, "EV-042, EV-043") {
		t.Errorf("EV refs joined missing: %q", text)
	}
}

func TestAppend_AllocatesNextIDFromHighWater(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.DomainDecisions)

	seed := "| ID | Date | Context | Decision | Rationale | Consequence | Supporting EV |\n" +
		"|----|------|---------|----------|-----------|-------------|---------------|\n" +
		"| DD-007 | 2026-04-01 | c | d | r | x | EV-001 |\n"
	if err := os.WriteFile(path, []byte(seed), 0o600); err != nil {
		t.Fatal(err)
	}

	id, err := decision.Append(path, decision.Row{
		Date: "2026-05-16", Context: "c", Decision: "d",
		Rationale: "r", Consequence: "x",
		SupportingEV: []string{"EV-002"},
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "DD-008" {
		t.Errorf("next id: want DD-008; got %q", id)
	}
}
