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
	expected := []string{
		"ctx-code-review", "ctx-status",
		"ctx-history", "ctx-brainstorm",
	}
	for _, exp := range expected {
		if !skillSet[exp] {
			t.Errorf("missing expected skill: %s", exp)
		}
	}
}

func TestSkillContent(t *testing.T) {
	content, err := FS.ReadFile(path.Join(
		asset.DirClaudeSkills,
		"ctx-history",
		asset.FileSKILLMd,
	))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "history") {
		t.Error("ctx-history SKILL.md does not contain 'history'")
	}
	if !strings.HasPrefix(string(content), "---") {
		t.Error("ctx-history SKILL.md missing frontmatter")
	}
}

func TestSkillReference(t *testing.T) {
	refPath := path.Join(
		asset.DirClaudeSkills, "ctx-skill-audit",
		asset.DirReferences,
		"anthropic-best-practices.md",
	)
	content, err := FS.ReadFile(refPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(content), "Anthropic") {
		t.Error("anthropic-best-practices.md does not contain 'Anthropic'")
	}
}

func TestListSkillReferences(t *testing.T) {
	refDir := path.Join(
		asset.DirClaudeSkills,
		"ctx-skill-audit",
		asset.DirReferences,
	)
	entries, err := FS.ReadDir(refDir)
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
	noRefDir := path.Join(
		asset.DirClaudeSkills,
		"ctx-status",
		asset.DirReferences,
	)
	_, err := FS.ReadDir(noRefDir)
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

func TestSchemaCoversCtxRC(t *testing.T) {
	// Parse the schema to get its property keys.
	schemaData, readErr := FS.ReadFile(asset.PathCtxrcSchema)
	if readErr != nil {
		t.Fatalf("read schema: %v", readErr)
	}
	var schema struct {
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if parseErr := json.Unmarshal(schemaData, &schema); parseErr != nil {
		t.Fatalf("parse schema: %v", parseErr)
	}

	// Parse a zero-value CtxRC to YAML then back to a map to get yaml tags.
	// We marshal a struct with all fields set to get every key emitted.
	type ctxRC struct {
		Profile             string `yaml:"profile"`
		TokenBudget         int    `yaml:"token_budget"`
		PriorityOrder       []int  `yaml:"priority_order"`
		AutoArchive         bool   `yaml:"auto_archive"`
		ArchiveAfterDays    int    `yaml:"archive_after_days"`
		ScratchpadEncrypt   *bool  `yaml:"scratchpad_encrypt"`
		EntryCountLearnings int    `yaml:"entry_count_learnings"`
		EntryCountDecisions int    `yaml:"entry_count_decisions"`
		ConventionLineCount int    `yaml:"convention_line_count"`
		InjectionTokenWarn  int    `yaml:"injection_token_warn"`
		ContextWindow       int    `yaml:"context_window"`
		BillingTokenWarn    int    `yaml:"billing_token_warn"`
		EventLog            bool   `yaml:"event_log"`
		KeyRotationDays     int    `yaml:"key_rotation_days"`
		TaskNudgeInterval   int    `yaml:"task_nudge_interval"`
		KeyPathOverride     string `yaml:"key_path"`
		StaleAgeDays        int    `yaml:"stale_age_days"`
		SessionPrefixes     []int  `yaml:"session_prefixes"`
		CompanionCheck      *bool  `yaml:"companion_check"`
		ClassifyRules       []int  `yaml:"classify_rules"`
		SpecSignalWords     []int  `yaml:"spec_signal_words"`
		SpecNudgeMinLen     int    `yaml:"spec_nudge_min_len"`
		Notify              *int   `yaml:"notify"`
		FreshnessFiles      []int  `yaml:"freshness_files"`
		Tool                string `yaml:"tool"`
		Steering            *int   `yaml:"steering"`
		Hooks               *int   `yaml:"hooks"`
		ProvenanceRequired  *int   `yaml:"provenance_required"`
	}
	yamlBytes, marshalErr := yaml.Marshal(ctxRC{})
	if marshalErr != nil {
		t.Fatalf("marshal: %v", marshalErr)
	}
	var structKeys map[string]any
	if unmarshalErr := yaml.Unmarshal(yamlBytes, &structKeys); unmarshalErr != nil {
		t.Fatalf("unmarshal: %v", unmarshalErr)
	}

	// Every struct field must appear in schema.
	for key := range structKeys {
		if _, ok := schema.Properties[key]; !ok {
			t.Errorf("CtxRC field %q has no schema property", key)
		}
	}
	// Every schema property must appear in struct.
	for key := range schema.Properties {
		if _, ok := structKeys[key]; !ok {
			t.Errorf("schema property %q has no CtxRC field", key)
		}
	}
}

func TestHookMessageRegistry(t *testing.T) {
	data, readErr := FS.ReadFile(asset.PathMessageRegistry)
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
	wantHooks := []string{
		"qa-reminder",
		"check-context-size",
		"block-non-path-ctx",
	}
	for _, exp := range wantHooks {
		if !hookSet[exp] {
			t.Errorf("missing expected hook: %s", exp)
		}
	}
}

func TestHookMessage_ReadVariant(t *testing.T) {
	gatePath := path.Join(
		asset.DirHooksMessages,
		"qa-reminder", "gate.txt",
	)
	content, readErr := FS.ReadFile(gatePath)
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
