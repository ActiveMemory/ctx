//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package steering

import (
	"testing"
)

func TestFilter_AlwaysIncludedRegardlessOfPrompt(t *testing.T) {
	files := []*SteeringFile{
		{Name: "always-on", Inclusion: InclusionAlways, Priority: 50},
	}

	got := Filter(files, "", nil, "")
	if len(got) != 1 || got[0].Name != "always-on" {
		t.Errorf("always-inclusion file should be included with empty prompt, got %v", names(got))
	}

	got = Filter(files, "completely unrelated prompt", nil, "")
	if len(got) != 1 || got[0].Name != "always-on" {
		t.Errorf("always-inclusion file should be included with any prompt, got %v", names(got))
	}
}

func TestFilter_AutoIncludedWhenPromptMatchesDescription(t *testing.T) {
	files := []*SteeringFile{
		{Name: "api-rules", Inclusion: InclusionAuto, Description: "REST API", Priority: 50},
	}

	got := Filter(files, "I need help with REST API design", nil, "")
	if len(got) != 1 || got[0].Name != "api-rules" {
		t.Errorf("auto file should match when prompt contains description, got %v", names(got))
	}

	// Case-insensitive match.
	got = Filter(files, "working on rest api endpoints", nil, "")
	if len(got) != 1 {
		t.Errorf("auto match should be case-insensitive, got %v", names(got))
	}
}

func TestFilter_AutoExcludedWhenPromptDoesNotMatch(t *testing.T) {
	files := []*SteeringFile{
		{Name: "api-rules", Inclusion: InclusionAuto, Description: "REST API", Priority: 50},
	}

	got := Filter(files, "fix the database migration", nil, "")
	if len(got) != 0 {
		t.Errorf("auto file should be excluded when prompt doesn't match, got %v", names(got))
	}
}

func TestFilter_ManualIncludedOnlyWhenNamed(t *testing.T) {
	files := []*SteeringFile{
		{Name: "security", Inclusion: InclusionManual, Priority: 50},
	}

	got := Filter(files, "anything", nil, "")
	if len(got) != 0 {
		t.Errorf("manual file should be excluded without explicit name, got %v", names(got))
	}

	got = Filter(files, "anything", []string{"security"}, "")
	if len(got) != 1 || got[0].Name != "security" {
		t.Errorf("manual file should be included when named, got %v", names(got))
	}
}

func TestFilter_PriorityOrdering(t *testing.T) {
	files := []*SteeringFile{
		{Name: "low", Inclusion: InclusionAlways, Priority: 90},
		{Name: "high", Inclusion: InclusionAlways, Priority: 10},
		{Name: "mid", Inclusion: InclusionAlways, Priority: 50},
	}

	got := Filter(files, "", nil, "")
	if len(got) != 3 {
		t.Fatalf("expected 3 files, got %d", len(got))
	}
	want := []string{"high", "mid", "low"}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("position %d: got %q, want %q", i, got[i].Name, name)
		}
	}
}

func TestFilter_AlphabeticalTieBreaking(t *testing.T) {
	files := []*SteeringFile{
		{Name: "charlie", Inclusion: InclusionAlways, Priority: 50},
		{Name: "alpha", Inclusion: InclusionAlways, Priority: 50},
		{Name: "bravo", Inclusion: InclusionAlways, Priority: 50},
	}

	got := Filter(files, "", nil, "")
	if len(got) != 3 {
		t.Fatalf("expected 3 files, got %d", len(got))
	}
	want := []string{"alpha", "bravo", "charlie"}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("position %d: got %q, want %q", i, got[i].Name, name)
		}
	}
}

func TestFilter_ToolFilterExcludesNonMatchingTool(t *testing.T) {
	files := []*SteeringFile{
		{Name: "cursor-only", Inclusion: InclusionAlways, Priority: 50, Tools: []string{"claude", "cursor"}},
	}

	got := Filter(files, "", nil, "kiro")
	if len(got) != 0 {
		t.Errorf("file with tools=[claude,cursor] should be excluded for tool=kiro, got %v", names(got))
	}
}

func TestFilter_EmptyToolsListIncludedForAnyTool(t *testing.T) {
	files := []*SteeringFile{
		{Name: "universal", Inclusion: InclusionAlways, Priority: 50, Tools: nil},
	}

	got := Filter(files, "", nil, "kiro")
	if len(got) != 1 || got[0].Name != "universal" {
		t.Errorf("file with empty tools should be included for any tool, got %v", names(got))
	}
}

func TestFilter_EmptyToolParameterSkipsToolFiltering(t *testing.T) {
	files := []*SteeringFile{
		{Name: "restricted", Inclusion: InclusionAlways, Priority: 50, Tools: []string{"cursor"}},
		{Name: "universal", Inclusion: InclusionAlways, Priority: 50, Tools: nil},
	}

	got := Filter(files, "", nil, "")
	if len(got) != 2 {
		t.Errorf("empty tool param should skip tool filtering, got %v", names(got))
	}
}

// names extracts file names for readable test output.
func names(files []*SteeringFile) []string {
	out := make([]string, len(files))
	for i, f := range files {
		out[i] = f.Name
	}
	return out
}
