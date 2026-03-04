//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"fmt"
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

	// Verify project-root and claude files are NOT in the list
	excluded := []string{
		"CLAUDE.md",
		"IMPLEMENTATION_PLAN.md",
		"Makefile.ctx",
	}
	for _, ex := range excluded {
		if templateSet[ex] {
			t.Errorf("List() should not contain %s (project-root file)", ex)
		}
	}
}

func TestClaudeMd(t *testing.T) {
	content, err := ClaudeMd()
	if err != nil {
		t.Fatalf("ClaudeMd() unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "Context") {
		t.Error("ClaudeMd() content does not contain \"Context\"")
	}
}

func TestProjectFile(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "IMPLEMENTATION_PLAN.md exists",
			file:        "IMPLEMENTATION_PLAN.md",
			wantContain: "Implementation",
			wantErr:     false,
		},
		{
			name:        "Makefile.ctx exists",
			file:        "Makefile.ctx",
			wantContain: "ctx",
			wantErr:     false,
		},
		{
			name:    "nonexistent project file returns error",
			file:    "NONEXISTENT.md",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := ProjectFile(tt.file)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ProjectFile(%q) expected error, got nil", tt.file)
				}
				return
			}
			if err != nil {
				t.Errorf("ProjectFile(%q) unexpected error: %v", tt.file, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("ProjectFile(%q) content does not contain %q", tt.file, tt.wantContain)
			}
		})
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

func TestListPromptTemplates(t *testing.T) {
	templates, err := ListPromptTemplates()
	if err != nil {
		t.Fatalf("ListPromptTemplates() unexpected error: %v", err)
	}

	if len(templates) == 0 {
		t.Error("ListPromptTemplates() returned empty list")
	}

	expected := []string{
		"code-review.md",
		"refactor.md",
		"explain.md",
	}

	templateSet := make(map[string]bool)
	for _, name := range templates {
		templateSet[name] = true
	}

	for _, exp := range expected {
		if !templateSet[exp] {
			t.Errorf("ListPromptTemplates() missing expected template: %s", exp)
		}
	}
}

func TestGetPromptTemplate(t *testing.T) {
	tests := []struct {
		name        string
		template    string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "code-review.md exists",
			template:    "code-review.md",
			wantContain: "Review",
			wantErr:     false,
		},
		{
			name:        "refactor.md exists",
			template:    "refactor.md",
			wantContain: "Refactor",
			wantErr:     false,
		},
		{
			name:        "explain.md exists",
			template:    "explain.md",
			wantContain: "Explain",
			wantErr:     false,
		},
		{
			name:     "nonexistent prompt template returns error",
			template: "nonexistent.md",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := PromptTemplate(tt.template)
			if tt.wantErr {
				if err == nil {
					t.Errorf("PromptTemplate(%q) expected error, got nil", tt.template)
				}
				return
			}
			if err != nil {
				t.Errorf("PromptTemplate(%q) unexpected error: %v", tt.template, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("PromptTemplate(%q) content does not contain %q", tt.template, tt.wantContain)
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
		"ctx-prompt",
		"ctx-status",
		"ctx-recall",
		"ctx-brainstorm",
		"ctx-check-links",
		"ctx-sanitize-permissions",
		"ctx-skill-audit",
		"ctx-skill-creator",
		"ctx-spec",
		"ctx-verify",
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

func TestSkillReference(t *testing.T) {
	content, err := SkillReference("ctx-skill-audit", "anthropic-best-practices.md")
	if err != nil {
		t.Fatalf("SkillReference() error: %v", err)
	}
	if !strings.Contains(string(content), "Anthropic") {
		t.Error("anthropic-best-practices.md does not contain 'Anthropic'")
	}
}

func TestListSkillReferences(t *testing.T) {
	refs, err := ListSkillReferences("ctx-skill-audit")
	if err != nil {
		t.Fatalf("ListSkillReferences() error: %v", err)
	}

	if len(refs) == 0 {
		t.Error("ListSkillReferences(ctx-skill-audit) returned empty list")
	}

	found := false
	for _, name := range refs {
		if name == "anthropic-best-practices.md" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ListSkillReferences() missing anthropic-best-practices.md")
	}
}

func TestListSkillReferencesNonexistent(t *testing.T) {
	_, err := ListSkillReferences("ctx-status")
	if err == nil {
		t.Error("ListSkillReferences(ctx-status) expected error for skill without references")
	}
}

func TestWhyDoc(t *testing.T) {
	tests := []struct {
		name        string
		doc         string
		wantContain string
		wantErr     bool
	}{
		{
			name:        "manifesto exists",
			doc:         "manifesto",
			wantContain: "Manifesto",
			wantErr:     false,
		},
		{
			name:        "about exists",
			doc:         "about",
			wantContain: "ctx",
			wantErr:     false,
		},
		{
			name:        "design-invariants exists",
			doc:         "design-invariants",
			wantContain: "Invariants",
			wantErr:     false,
		},
		{
			name:    "nonexistent doc returns error",
			doc:     "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := WhyDoc(tt.doc)
			if tt.wantErr {
				if err == nil {
					t.Errorf("WhyDoc(%q) expected error, got nil", tt.doc)
				}
				return
			}
			if err != nil {
				t.Errorf("WhyDoc(%q) unexpected error: %v", tt.doc, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("WhyDoc(%q) content does not contain %q", tt.doc, tt.wantContain)
			}
		})
	}
}

func TestListWhyDocs(t *testing.T) {
	docs, err := ListWhyDocs()
	if err != nil {
		t.Fatalf("ListWhyDocs() unexpected error: %v", err)
	}

	expected := []string{"about", "design-invariants", "manifesto"}

	docSet := make(map[string]bool)
	for _, name := range docs {
		docSet[name] = true
	}

	for _, exp := range expected {
		if !docSet[exp] {
			t.Errorf("ListWhyDocs() missing expected doc: %s", exp)
		}
	}

	if len(docs) != len(expected) {
		t.Errorf("ListWhyDocs() returned %d docs, expected %d", len(docs), len(expected))
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

func TestSchema(t *testing.T) {
	data, err := Schema()
	if err != nil {
		t.Fatalf("Schema() unexpected error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "$schema") {
		t.Error("Schema() content does not contain $schema")
	}
	if !strings.Contains(content, "additionalProperties") {
		t.Error("Schema() content does not contain additionalProperties")
	}
	if !strings.Contains(content, "ctx.ist") {
		t.Error("Schema() content does not contain ctx.ist $id")
	}
}

func TestDefaultPermissions(t *testing.T) {
	t.Run("allow list is non-empty", func(t *testing.T) {
		allow := DefaultAllowPermissions()
		if len(allow) == 0 {
			t.Fatal("DefaultAllowPermissions() returned empty list")
		}
	})

	t.Run("deny list is non-empty", func(t *testing.T) {
		deny := DefaultDenyPermissions()
		if len(deny) == 0 {
			t.Fatal("DefaultDenyPermissions() returned empty list")
		}
	})

	t.Run("allow list contains Bash(ctx:*)", func(t *testing.T) {
		allowSet := make(map[string]bool)
		for _, p := range DefaultAllowPermissions() {
			allowSet[p] = true
		}
		if !allowSet["Bash(ctx:*)"] {
			t.Error("allow list missing: Bash(ctx:*)")
		}
	})

	t.Run("allow list covers all bundled skills", func(t *testing.T) {
		skills, listErr := ListSkills()
		if listErr != nil {
			t.Fatalf("ListSkills() error: %v", listErr)
		}

		allowSet := make(map[string]bool)
		for _, p := range DefaultAllowPermissions() {
			allowSet[p] = true
		}

		for _, skill := range skills {
			entry := fmt.Sprintf("Skill(%s)", skill)
			if !allowSet[entry] {
				t.Errorf("allow list missing bundled skill: %s", entry)
			}
		}
	})

	t.Run("deny list contains dangerous patterns", func(t *testing.T) {
		denySet := make(map[string]bool)
		for _, p := range DefaultDenyPermissions() {
			denySet[p] = true
		}

		expected := []string{
			"Bash(sudo *)",
			"Bash(git push *)",
			"Bash(rm -rf /*)",
			"Bash(curl *)",
			"Read(**/.env)",
			"Edit(**/.env)",
		}
		for _, e := range expected {
			if !denySet[e] {
				t.Errorf("deny list missing: %s", e)
			}
		}
	})
}
