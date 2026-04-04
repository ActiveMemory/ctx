//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package skill

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSkillManifest(t *testing.T, dir, name, content string) {
	t.Helper()
	skillDir := filepath.Join(dir, name)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, skillManifest), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadAll_NonExistentDir(t *testing.T) {
	skills, err := LoadAll("/nonexistent/path")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(skills) != 0 {
		t.Fatalf("expected empty slice, got %d skills", len(skills))
	}
}

func TestLoadAll_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	skills, err := LoadAll(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(skills) != 0 {
		t.Fatalf("expected empty slice, got %d skills", len(skills))
	}
}

func TestLoadAll_MultipleSkills(t *testing.T) {
	dir := t.TempDir()

	writeSkillManifest(t, dir, "alpha", `---
name: alpha
description: Alpha skill
---
Alpha body content
`)
	writeSkillManifest(t, dir, "beta", `---
name: beta
description: Beta skill
---
Beta body content
`)

	skills, err := LoadAll(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 2 {
		t.Fatalf("expected 2 skills, got %d", len(skills))
	}

	// os.ReadDir returns entries sorted by name.
	if skills[0].Name != "alpha" {
		t.Errorf("expected first skill name 'alpha', got %q", skills[0].Name)
	}
	if skills[1].Name != "beta" {
		t.Errorf("expected second skill name 'beta', got %q", skills[1].Name)
	}
}

func TestLoadAll_SkipsNonDirectories(t *testing.T) {
	dir := t.TempDir()

	writeSkillManifest(t, dir, "valid", `---
name: valid
description: A valid skill
---
Body
`)
	// Create a regular file at the top level — should be skipped.
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("ignore"), 0o644); err != nil {
		t.Fatal(err)
	}

	skills, err := LoadAll(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}
}

func TestLoad_ValidSkill(t *testing.T) {
	dir := t.TempDir()
	writeSkillManifest(t, dir, "react-patterns", `---
name: react-patterns
description: React component patterns
---
# React Patterns
- Use functional components
`)

	sk, err := Load(dir, "react-patterns")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sk.Name != "react-patterns" {
		t.Errorf("expected name 'react-patterns', got %q", sk.Name)
	}
	if sk.Description != "React component patterns" {
		t.Errorf("expected description 'React component patterns', got %q", sk.Description)
	}
	if sk.Body != "# React Patterns\n- Use functional components\n" {
		t.Errorf("unexpected body: %q", sk.Body)
	}
	expectedDir := filepath.Join(dir, "react-patterns")
	if sk.Dir != expectedDir {
		t.Errorf("expected dir %q, got %q", expectedDir, sk.Dir)
	}
}

func TestLoad_MissingManifest(t *testing.T) {
	dir := t.TempDir()
	// Create subdirectory without SKILL.md.
	if err := os.MkdirAll(filepath.Join(dir, "empty-skill"), 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir, "empty-skill")
	if err == nil {
		t.Fatal("expected error for missing SKILL.md")
	}
}

func TestLoad_InvalidFrontmatter(t *testing.T) {
	dir := t.TempDir()
	writeSkillManifest(t, dir, "bad", `---
name: [invalid yaml
---
Body
`)

	_, err := Load(dir, "bad")
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_MissingDelimiters(t *testing.T) {
	dir := t.TempDir()
	writeSkillManifest(t, dir, "no-fm", `Just plain markdown without frontmatter`)

	_, err := Load(dir, "no-fm")
	if err == nil {
		t.Fatal("expected error for missing frontmatter delimiters")
	}
}

func TestLoad_EmptyDescription(t *testing.T) {
	dir := t.TempDir()
	writeSkillManifest(t, dir, "minimal", `---
name: minimal
---
Minimal body
`)

	sk, err := Load(dir, "minimal")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sk.Name != "minimal" {
		t.Errorf("expected name 'minimal', got %q", sk.Name)
	}
	if sk.Description != "" {
		t.Errorf("expected empty description, got %q", sk.Description)
	}
}
