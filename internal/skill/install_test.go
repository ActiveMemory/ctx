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

	cfgSkill "github.com/ActiveMemory/ctx/internal/config/skill"
)

func TestInstall_ValidSkill(t *testing.T) {
	source := t.TempDir()
	skillsDir := t.TempDir()

	manifest := `---
name: test-skill
description: A test skill
---
# Instructions
Do the thing.
`
	if err := os.WriteFile(filepath.Join(source, cfgSkill.SkillManifest), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	// Add an extra file to verify full directory copy.
	if err := os.WriteFile(filepath.Join(source, "extra.md"), []byte("extra"), 0o644); err != nil {
		t.Fatal(err)
	}

	sk, err := Install(source, skillsDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if sk.Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got %q", sk.Name)
	}
	if sk.Description != "A test skill" {
		t.Errorf("expected description 'A test skill', got %q", sk.Description)
	}
	if sk.Dir != filepath.Join(skillsDir, "test-skill") {
		t.Errorf("expected dir %q, got %q", filepath.Join(skillsDir, "test-skill"), sk.Dir)
	}

	// Verify SKILL.md was copied.
	if _, err := os.Stat(filepath.Join(skillsDir, "test-skill", cfgSkill.SkillManifest)); err != nil {
		t.Errorf("SKILL.md not copied: %v", err)
	}
	// Verify extra file was copied.
	if _, err := os.Stat(filepath.Join(skillsDir, "test-skill", "extra.md")); err != nil {
		t.Errorf("extra.md not copied: %v", err)
	}
}

func TestInstall_CopiesSubdirectories(t *testing.T) {
	source := t.TempDir()
	skillsDir := t.TempDir()

	manifest := `---
name: nested-skill
---
Body
`
	if err := os.WriteFile(filepath.Join(source, cfgSkill.SkillManifest), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(source, "templates")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "tmpl.txt"), []byte("template"), 0o644); err != nil {
		t.Fatal(err)
	}

	sk, err := Install(source, skillsDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	copied := filepath.Join(sk.Dir, "templates", "tmpl.txt")
	data, err := os.ReadFile(copied)
	if err != nil {
		t.Fatalf("subdirectory file not copied: %v", err)
	}
	if string(data) != "template" {
		t.Errorf("expected 'template', got %q", string(data))
	}
}

func TestInstall_MissingManifest(t *testing.T) {
	source := t.TempDir()
	skillsDir := t.TempDir()

	_, err := Install(source, skillsDir)
	if err == nil {
		t.Fatal("expected error for missing SKILL.md")
	}
}

func TestInstall_InvalidFrontmatter(t *testing.T) {
	source := t.TempDir()
	skillsDir := t.TempDir()

	manifest := `---
name: [broken yaml
---
Body
`
	if err := os.WriteFile(filepath.Join(source, cfgSkill.SkillManifest), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Install(source, skillsDir)
	if err == nil {
		t.Fatal("expected error for invalid YAML frontmatter")
	}
}

func TestInstall_MissingName(t *testing.T) {
	source := t.TempDir()
	skillsDir := t.TempDir()

	manifest := `---
description: no name field
---
Body
`
	if err := os.WriteFile(filepath.Join(source, cfgSkill.SkillManifest), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Install(source, skillsDir)
	if err == nil {
		t.Fatal("expected error for missing name field")
	}
}

func TestInstall_NoFrontmatterDelimiters(t *testing.T) {
	source := t.TempDir()
	skillsDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(source, cfgSkill.SkillManifest), []byte("just plain text"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Install(source, skillsDir)
	if err == nil {
		t.Fatal("expected error for missing frontmatter delimiters")
	}
}
