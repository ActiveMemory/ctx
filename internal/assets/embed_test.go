//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"strings"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "CONSTITUTION.md exists",
			template:    "CONSTITUTION.md",
			wantContain: "Constitution",
			wantErr:     false,
		},
		{
			name:        "TASKS.md exists",
			template:    "TASKS.md",
			wantContain: "Tasks",
			wantErr:     false,
		},
		{
			name:        "DECISIONS.md exists",
			template:    "DECISIONS.md",
			wantContain: "Decisions",
			wantErr:     false,
		},
		{
			name:        "LEARNINGS.md exists",
			template:    "LEARNINGS.md",
			wantContain: "Learnings",
			wantErr:     false,
		},
		{
			name:        "CONVENTIONS.md exists",
			template:    "CONVENTIONS.md",
			wantContain: "Conventions",
			wantErr:     false,
		},
		{
			name:        "ARCHITECTURE.md exists",
			template:    "ARCHITECTURE.md",
			wantContain: "Architecture",
			wantErr:     false,
		},
		{
			name:        "AGENT_PLAYBOOK.md exists",
			template:    "AGENT_PLAYBOOK.md",
			wantContain: "Agent Playbook",
			wantErr:     false,
		},
		{
			name:        "GLOSSARY.md exists",
			template:    "GLOSSARY.md",
			wantContain: "Glossary",
			wantErr:     false,
		},
		{
			name:        "CLAUDE.md exists",
			template:    "CLAUDE.md",
			wantContain: "Context",
			wantErr:     false,
		},
		{
			name:     "nonexistent template returns error",
			template: "NONEXISTENT.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := Template(tt.template)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Template(%q) expected error, got nil", tt.template)
				}
				return
			}
			if err != nil {
				t.Errorf("Template(%q) unexpected error: %v", tt.template, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("Template(%q) content does not contain %q", tt.template, tt.wantContain)
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	templates, err := List()
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	if len(templates) == 0 {
		t.Error("List() returned empty list")
	}

	// Check for required templates
	required := []string{
		"CONSTITUTION.md",
		"TASKS.md",
		"DECISIONS.md",
		"LEARNINGS.md",
	}

	templateSet := make(map[string]bool)
	for _, name := range templates {
		templateSet[name] = true
	}

	for _, req := range required {
		if !templateSet[req] {
			t.Errorf("List() missing required template: %s", req)
		}
	}
}

func TestListEntryTemplates(t *testing.T) {
	templates, err := ListEntry()
	if err != nil {
		t.Fatalf("ListEntry() unexpected error: %v", err)
	}

	if len(templates) == 0 {
		t.Error("ListEntry() returned empty list")
	}

	// Check for expected entry templates
	expected := []string{
		"learning.md",
		"decision.md",
	}

	templateSet := make(map[string]bool)
	for _, name := range templates {
		templateSet[name] = true
	}

	for _, exp := range expected {
		if !templateSet[exp] {
			t.Errorf("ListEntry() missing expected template: %s", exp)
		}
	}
}

func TestGetEntryTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "learning.md exists",
			template:    "learning.md",
			wantContain: "Context",
			wantErr:     false,
		},
		{
			name:        "decision.md exists",
			template:    "decision.md",
			wantContain: "Context",
			wantErr:     false,
		},
		{
			name:     "nonexistent entry template returns error",
			template: "nonexistent.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := Entry(tt.template)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Entry(%q) expected error, got nil", tt.template)
				}
				return
			}
			if err != nil {
				t.Errorf("Entry(%q) unexpected error: %v", tt.template, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("Entry(%q) content does not contain %q", tt.template, tt.wantContain)
			}
		})
	}
}

func TestListSkills(t *testing.T) {
	skills, err := ListSkills()
	if err != nil {
		t.Fatalf("ListSkills() unexpected error: %v", err)
	}

	if len(skills) == 0 {
		t.Error("ListSkills() returned empty list")
	}

	// Check for expected skills (directory names, not files)
	expected := []string{
		"ctx-status",
		"ctx-recall",
	}

	skillSet := make(map[string]bool)
	for _, name := range skills {
		skillSet[name] = true
	}

	for _, exp := range expected {
		if !skillSet[exp] {
			t.Errorf("ListSkills() missing expected skill: %s", exp)
		}
	}
}

func TestSkillContent(t *testing.T) {
	content, err := SkillContent("ctx-recall")
	if err != nil {
		t.Fatalf("SkillContent(ctx-recall) error: %v", err)
	}
	if !strings.Contains(string(content), "recall") {
		t.Error("ctx-recall SKILL.md does not contain 'recall'")
	}
	// Verify it's a valid SKILL.md with frontmatter
	if !strings.HasPrefix(string(content), "---") {
		t.Error("ctx-recall SKILL.md missing frontmatter")
	}
}

func TestListTools(t *testing.T) {
	tools, err := ListTools()
	if err != nil {
		t.Fatalf("ListTools() unexpected error: %v", err)
	}

	if len(tools) == 0 {
		t.Error("ListTools() returned empty list")
	}

	toolSet := make(map[string]bool)
	for _, name := range tools {
		toolSet[name] = true
	}

	if !toolSet["context-watch.sh"] {
		t.Error("ListTools() missing expected tool: context-watch.sh")
	}
}

func TestToolContent(t *testing.T) {
	content, err := Tool("context-watch.sh")
	if err != nil {
		t.Fatalf("Tool(context-watch.sh) error: %v", err)
	}
	if !strings.Contains(string(content), "Context Monitor") {
		t.Error("context-watch.sh does not contain 'Context Monitor'")
	}
	if !strings.HasPrefix(string(content), "#!/bin/bash") {
		t.Error("context-watch.sh missing bash shebang")
	}
}

func TestPluginVersion(t *testing.T) {
	ver, err := PluginVersion()
	if err != nil {
		t.Fatalf("PluginVersion() error: %v", err)
	}
	if ver == "" {
		t.Error("PluginVersion() returned empty string")
	}
	// Should be a semver-like string
	if !strings.Contains(ver, ".") {
		t.Errorf("PluginVersion() = %q, expected semver format", ver)
	}
}
