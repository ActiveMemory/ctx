//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package consolidate

import (
	"strings"
	"testing"
)

func TestConsolidateToolRuns(t *testing.T) {
	content := strings.Join([]string{
		"### 1. Assistant (10:00:00)",
		"",
		"Read file.go",
		"",
		"### 2. Assistant (10:00:01)",
		"",
		"Read file.go",
		"",
		"### 3. Assistant (10:00:02)",
		"",
		"Read file.go",
		"",
		"### 4. User (10:00:03)",
		"",
		"Done",
	}, "\n")

	got := ToolRuns(content)

	if !strings.Contains(got, "\u00d73)") {
		t.Errorf("expected (x3) count marker, got:\n%s", got)
	}
	if strings.Contains(got, "### 2. Assistant") {
		t.Error("duplicate turns should be collapsed")
	}
	if !strings.Contains(got, "### 4. User") {
		t.Error("different turn should be preserved")
	}
}

func TestConsolidateToolRuns_DifferentTools(t *testing.T) {
	content := strings.Join([]string{
		"### 1. Assistant (10:00:00)",
		"",
		"Read file.go",
		"",
		"### 2. Assistant (10:00:01)",
		"",
		"Write file.go",
	}, "\n")

	got := ToolRuns(content)

	if strings.Contains(got, "\u00d7") {
		t.Error("different bodies should not be consolidated")
	}
	if !strings.Contains(got, "### 1. Assistant") {
		t.Error("first turn missing")
	}
	if !strings.Contains(got, "### 2. Assistant") {
		t.Error("second turn missing")
	}
}

func TestConsolidateToolRuns_SingleTurn(t *testing.T) {
	content := "### 1. Assistant (10:00:00)\n\nSingle turn\n"
	got := ToolRuns(content)

	if strings.Contains(got, "\u00d7") {
		t.Error("single turn should not have count marker")
	}
}
