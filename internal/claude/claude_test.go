//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package claude

import (
	"embed"
	"errors"
	"strings"
	"testing"

	"github.com/ActiveMemory/ctx/internal/assets"
)

func TestSkills(t *testing.T) {
	skills, err := Skills()
	if err != nil {
		t.Fatalf("Skills() unexpected error: %v", err)
	}

	if len(skills) == 0 {
		t.Error("Skills() returned empty list")
	}

	// Check that all entries are skill directory names (no extension)
	for _, skill := range skills {
		if strings.Contains(skill, ".") {
			t.Errorf("Skills() returned name with extension: %s", skill)
		}
	}
}

func TestSkillContent(t *testing.T) {
	// First get the list of skills to test with
	skills, err := Skills()
	if err != nil {
		t.Fatalf("Skills() failed: %v", err)
	}

	if len(skills) == 0 {
		t.Skip("no skills available to test")
	}

	// Test getting the first skill
	content, err := SkillContent(skills[0])
	if err != nil {
		t.Errorf("SkillContent(%q) unexpected error: %v", skills[0], err)
	}
	if len(content) == 0 {
		t.Errorf("SkillContent(%q) returned empty content", skills[0])
	}

	// Verify it's a valid SKILL.md with frontmatter
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---") {
		t.Errorf("SkillContent(%q) missing frontmatter", skills[0])
	}

	// Test getting nonexistent skill
	_, err = SkillContent("nonexistent-skill")
	if err == nil {
		t.Error("SkillContent(nonexistent) expected error, got nil")
	}
}

func TestSkillContentAllSkills(t *testing.T) {
	skills, err := Skills()
	if err != nil {
		t.Fatalf("Skills() failed: %v", err)
	}
	for _, name := range skills {
		content, err := SkillContent(name)
		if err != nil {
			t.Errorf("SkillContent(%q) error: %v", name, err)
			continue
		}
		if len(content) == 0 {
			t.Errorf("SkillContent(%q) returned empty content", name)
		}
	}
}

func TestSettingsStructure(t *testing.T) {
	settings := Settings{
		Permissions: PermissionsConfig{
			Allow: []string{"Bash(ctx status:*)", "Bash(ctx agent:*)"},
		},
	}

	if len(settings.Permissions.Allow) == 0 {
		t.Error("Settings.Permissions.Allow should not be empty")
	}
}

func TestErrSkillList(t *testing.T) {
	cause := errors.New("read dir failed")
	err := errSkillList(cause)
	if err == nil {
		t.Fatal("errSkillList() returned nil")
	}
	if !strings.Contains(err.Error(), "failed to list skills") {
		t.Errorf("errSkillList() error missing prefix: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("errSkillList() does not wrap the cause error")
	}
}

func TestErrSkillRead(t *testing.T) {
	cause := errors.New("not found")
	err := errSkillRead("my-skill", cause)
	if err == nil {
		t.Fatal("errSkillRead() returned nil")
	}
	if !strings.Contains(err.Error(), "my-skill") {
		t.Errorf("errSkillRead() error missing skill name: %v", err)
	}
	if !errors.Is(err, cause) {
		t.Error("errSkillRead() does not wrap the cause error")
	}
}

// TestScriptErrorPaths swaps assets.FS with an empty embed.FS to trigger
// error branches in skill functions.
func TestScriptErrorPaths(t *testing.T) {
	orig := assets.FS
	defer func() { assets.FS = orig }()
	assets.FS = embed.FS{} // empty FS causes all reads to fail

	if _, err := Skills(); err == nil {
		t.Error("Skills() expected error with empty FS")
	}
}
