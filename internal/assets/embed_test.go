//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package assets

import (
	"encoding/json"
	"path"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/config/asset"
	"github.com/ActiveMemory/ctx/internal/config/file"

	"gopkg.in/yaml.v3"
)

// TestDescKeysResolve lives in read/desc/desc_test.go where it can
// call lookup.Init() without an import cycle.

// TestDefaultPermissions lives in read/lookup/perm_test.go where it can
// call Init() without an import cycle.

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

	for _, req := range []string{"CONSTITUTION.md", "TASKS.md", "DECISIONS.md", "LEARNINGS.md"} {
		if !templateSet[req] {
			t.Errorf("missing required template: %s", req)
		}
	}
	for _, ex := range []string{"CLAUDE.md", "IMPLEMENTATION_PLAN.md", "Makefile.ctx"} {
		if templateSet[ex] {
			t.Errorf("should not contain project-root file: %s", ex)
		}
	}
}

func TestClaudeMd(t *testing.T) {
	content, err := FS.ReadFile(asset.PathCLAUDEMd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "Context") {
		t.Error("CLAUDE.md does not contain 'Context'")
	}
}

func TestProjectFile(t *testing.T) {
	tests := []struct {
		name        string
		file        string
		wantContain string
		wantErr     bool
	}{
		{"IMPLEMENTATION_PLAN.md exists", "IMPLEMENTATION_PLAN.md", "Implementation", false},
		{"Makefile.ctx exists", "Makefile.ctx", "ctx", false},
		{"nonexistent returns error", "NONEXISTENT.md", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := FS.ReadFile(path.Join(asset.DirProject, tt.file))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.file)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.file, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("content of %q does not contain %q", tt.file, tt.wantContain)
			}
		})
	}
}

func TestListPromptTemplates(t *testing.T) {
	entries, err := FS.ReadDir(asset.DirPromptTemplates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Error("returned empty list")
	}

	nameSet := make(map[string]bool)
	for _, e := range entries {
		nameSet[e.Name()] = true
	}
	for _, exp := range []string{"code-review.md", "refactor.md", "explain.md"} {
		if !nameSet[exp] {
			t.Errorf("missing expected template: %s", exp)
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
		{"code-review.md exists", "code-review.md", "Review", false},
		{"refactor.md exists", "refactor.md", "Refactor", false},
		{"explain.md exists", "explain.md", "Explain", false},
		{"nonexistent returns error", "nonexistent.md", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := FS.ReadFile(path.Join(asset.DirPromptTemplates, tt.template))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.template)
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

func TestListSkills(t *testing.T) {
	entries, err := FS.ReadDir(asset.DirClaudeSkills)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Error("returned empty list")
	}

	skillSet := make(map[string]bool)
	for _, e := range entries {
		if e.IsDir() {
			skillSet[e.Name()] = true
		}
	}
	for _, exp := range []string{"ctx-prompt", "ctx-status", "ctx-recall", "ctx-brainstorm"} {
		if !skillSet[exp] {
			t.Errorf("missing expected skill: %s", exp)
		}
	}
}

func TestSkillContent(t *testing.T) {
	content, err := FS.ReadFile(path.Join(asset.DirClaudeSkills, "ctx-recall", asset.FileSKILLMd))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "recall") {
		t.Error("ctx-recall SKILL.md does not contain 'recall'")
	}
	if !strings.HasPrefix(string(content), "---") {
		t.Error("ctx-recall SKILL.md missing frontmatter")
	}
}

func TestSkillReference(t *testing.T) {
	content, err := FS.ReadFile(path.Join(
		asset.DirClaudeSkills, "ctx-skill-audit", asset.DirReferences, "anthropic-best-practices.md",
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "Anthropic") {
		t.Error("anthropic-best-practices.md does not contain 'Anthropic'")
	}
}

func TestListSkillReferences(t *testing.T) {
	entries, err := FS.ReadDir(path.Join(asset.DirClaudeSkills, "ctx-skill-audit", asset.DirReferences))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) == 0 {
		t.Error("returned empty list")
	}

	found := false
	for _, e := range entries {
		if e.Name() == "anthropic-best-practices.md" {
			found = true
			break
		}
	}
	if !found {
		t.Error("missing anthropic-best-practices.md")
	}
}

func TestListSkillReferencesNonexistent(t *testing.T) {
	_, err := FS.ReadDir(path.Join(asset.DirClaudeSkills, "ctx-status", asset.DirReferences))
	if err == nil {
		t.Error("expected error for skill without references")
	}
}

func TestWhyDoc(t *testing.T) {
	tests := []struct {
		name        string
		doc         string
		wantContain string
		wantErr     bool
	}{
		{"manifesto exists", "manifesto", "Manifesto", false},
		{"about exists", "about", "ctx", false},
		{"design-invariants exists", "design-invariants", "Invariants", false},
		{"nonexistent returns error", "nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := FS.ReadFile(path.Join(asset.DirWhy, tt.doc+file.ExtMarkdown))
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q", tt.doc)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error for %q: %v", tt.doc, err)
				return
			}
			if !strings.Contains(string(content), tt.wantContain) {
				t.Errorf("content of %q does not contain %q", tt.doc, tt.wantContain)
			}
		})
	}
}

func TestListWhyDocs(t *testing.T) {
	entries, err := FS.ReadDir(asset.DirWhy)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"about", "design-invariants", "manifesto"}
	docSet := make(map[string]bool)
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, file.ExtMarkdown) {
			docSet[strings.TrimSuffix(name, file.ExtMarkdown)] = true
		}
	}

	for _, exp := range expected {
		if !docSet[exp] {
			t.Errorf("missing expected doc: %s", exp)
		}
	}
}

func TestPluginVersion(t *testing.T) {
	data, err := FS.ReadFile(asset.PathPluginJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var manifest map[string]json.RawMessage
	if unmarshalErr := json.Unmarshal(data, &manifest); unmarshalErr != nil {
		t.Fatalf("parse error: %v", unmarshalErr)
	}
	raw, ok := manifest[asset.JSONKeyVersion]
	if !ok {
		t.Fatal("plugin.json missing 'version' key")
	}
	var ver string
	if parseErr := json.Unmarshal(raw, &ver); parseErr != nil {
		t.Fatalf("version parse error: %v", parseErr)
	}
	if ver == "" {
		t.Error("version is empty")
	}
	if !strings.Contains(ver, ".") {
		t.Errorf("version = %q, expected semver format", ver)
	}
}

func TestSchema(t *testing.T) {
	data, err := FS.ReadFile(asset.PathCtxrcSchema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "$schema") {
		t.Error("does not contain $schema")
	}
	if !strings.Contains(content, "ctx.ist") {
		t.Error("does not contain ctx.ist $id")
	}
}

func TestHookMessageRegistry(t *testing.T) {
	data, readErr := FS.ReadFile(asset.PathHookRegistry)
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(data) == 0 {
		t.Fatal("returned empty data")
	}

	var entries []map[string]any
	if parseErr := yaml.Unmarshal(data, &entries); parseErr != nil {
		t.Fatalf("invalid YAML: %v", parseErr)
	}
	for i, entry := range entries {
		if _, ok := entry["hook"]; !ok {
			t.Errorf("entry %d missing 'hook' key", i)
		}
		if _, ok := entry["variant"]; !ok {
			t.Errorf("entry %d missing 'variant' key", i)
		}
	}
}

func TestListHookMessages(t *testing.T) {
	entries, listErr := FS.ReadDir(asset.DirHooksMessages)
	if listErr != nil {
		t.Fatalf("unexpected error: %v", listErr)
	}
	if len(entries) == 0 {
		t.Fatal("returned empty list")
	}

	hookSet := make(map[string]bool)
	for _, h := range entries {
		if h.IsDir() {
			hookSet[h.Name()] = true
		}
	}
	for _, exp := range []string{"qa-reminder", "check-context-size", "block-dangerous-commands"} {
		if !hookSet[exp] {
			t.Errorf("missing expected hook: %s", exp)
		}
	}
}

func TestHookMessage_ReadVariant(t *testing.T) {
	content, readErr := FS.ReadFile(path.Join(asset.DirHooksMessages, "qa-reminder", "gate.txt"))
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(content) == 0 {
		t.Fatal("returned empty content")
	}
}

func TestRalphTemplate(t *testing.T) {
	content, readErr := FS.ReadFile(path.Join(asset.DirRalph, "PROMPT.md"))
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(content) == 0 {
		t.Fatal("returned empty content")
	}
}

func TestMakefileCtx(t *testing.T) {
	content, readErr := FS.ReadFile(asset.PathMakefileCtx)
	if readErr != nil {
		t.Fatalf("unexpected error: %v", readErr)
	}
	if len(content) == 0 {
		t.Fatal("returned empty content")
	}
	if !strings.Contains(string(content), "ctx") {
		t.Error("content does not contain 'ctx'")
	}
}
