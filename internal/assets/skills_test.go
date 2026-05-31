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
