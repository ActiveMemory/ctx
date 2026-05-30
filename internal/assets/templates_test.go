//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"path"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
)

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantContain string
		wantErr     bool
	}{
		{"CONSTITUTION.md exists", "CONSTITUTION.md", "Constitution", false},
		{"TASKS.md exists", "TASKS.md", "Tasks", false},
		{"DECISIONS.md exists", "DECISIONS.md", "Decisions", false},
		{"LEARNINGS.md exists", "LEARNINGS.md", "Learnings", false},
		{"CONVENTIONS.md exists", "CONVENTIONS.md", "Conventions", false},
		{"ARCHITECTURE.md exists", "ARCHITECTURE.md", "Architecture", false},
		{"AGENT_PLAYBOOK.md exists", "AGENT_PLAYBOOK.md", "Agent Playbook", false},
		{"AGENT_PLAYBOOK_GATE.md exists", "AGENT_PLAYBOOK_GATE.md", "Agent Playbook (Gate)", false},
		{"GLOSSARY.md exists", "GLOSSARY.md", "Glossary", false},
		{"nonexistent template returns error", "NONEXISTENT.md", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := FS.ReadFile(path.Join(asset.DirContext, tt.template))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.template)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.template, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("content of %q does not contain %q", tt.template, tt.wantContain)
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	entries, err := FS.ReadDir(asset.DirContext)
	if err != nil {
		t.Fatalf("ReadDir() unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Error("ReadDir() returned empty list")
	}

	templateSet := make(map[string]bool)
	for _, e := range entries {
		templateSet[e.Name()] = true
	}

	required := []string{
		"CONSTITUTION.md", "TASKS.md",
		"DECISIONS.md", "LEARNINGS.md",
	}
	for _, req := range required {
		if !templateSet[req] {
			t.Errorf("missing required template: %s", req)
		}
	}
	for _, ex := range []string{"CLAUDE.md", "Makefile.ctx"} {
		if templateSet[ex] {
			t.Errorf("should not contain project-root file: %s", ex)
		}
	}
}
