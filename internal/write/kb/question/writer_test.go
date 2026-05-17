//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package question_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfgKB "github.com/ActiveMemory/ctx/internal/config/kb"
	"github.com/ActiveMemory/ctx/internal/write/kb/question"
)

func TestAppend_FirstRowGetsQ001(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.OutstandingQuestions)

	id, err := question.Append(path, question.Row{
		Question:                 "Does the widget folder permit nested skills?",
		WhyItMatters:             "Blocks widget-composition promotion to high.",
		WhatEvidenceWouldResolve: "An upstream spec statement.",
		OpenedAt:                 "2026-05-16",
		RelatedEV:                []string{"EV-042", "EV-043"},
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "Q-001" {
		t.Errorf("first id: want Q-001; got %q", id)
	}

	got, _ := os.ReadFile(path)
	text := string(got)
	if !strings.Contains(text, "| ID | Question |") {
		t.Errorf("header missing: %q", text)
	}
	if !strings.Contains(text, "Q-001") {
		t.Errorf("ID missing: %q", text)
	}
	if !strings.Contains(text, "EV-042, EV-043") {
		t.Errorf("related EV joined missing: %q", text)
	}
}

func TestAppend_AllocatesNextIDFromHighWater(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, cfgKB.OutstandingQuestions)

	seed := "| ID | Question | Why it matters | What evidence would resolve | Opened at | Related EV |\n" +
		"|----|----------|----------------|-----------------------------|-----------|------------|\n" +
		"| Q-003 | q | y | w | 2026-05-01 |  |\n" +
		"| Q-009 | q | y | w | 2026-05-10 |  |\n"
	if err := os.WriteFile(path, []byte(seed), 0o600); err != nil {
		t.Fatal(err)
	}

	id, err := question.Append(path, question.Row{
		Question: "q", WhyItMatters: "y",
		WhatEvidenceWouldResolve: "w", OpenedAt: "2026-05-16",
	})
	if err != nil {
		t.Fatalf("Append: %v", err)
	}
	if id != "Q-010" {
		t.Errorf("next id: want Q-010; got %q", id)
	}
}
